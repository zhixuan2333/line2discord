package main

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	log "github.com/sirupsen/logrus"
)

// WebHook
func WebHook(w http.ResponseWriter, req *http.Request) {
	events, err := LineBot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		id := getDiscordID(event)

		switch event.Type {
		case linebot.EventTypeMessage:

			switch message := event.Message.(type) {

			// Text message
			case *linebot.TextMessage:

				DiscordSendMessage(event.Source.UserID, id, message.Text)

			// Image message
			case *linebot.ImageMessage:

				cw, err := LineBot.GetMessageContent(message.ID).Do()
				if err != nil {
					log.Error("Get line file content", err)
					return
				}

				log.Info(cw.ContentType)

				defer cw.Content.Close()

				DiscordSendFile(event.Source.UserID, message.ID, id, cw.ContentType, cw.Content)

			// Video message
			case *linebot.VideoMessage:

				cw, err := LineBot.GetMessageContent(message.ID).Do()
				if err != nil {
					log.Error("Get line file content", err)
					return
				}
				DiscordSendFile(event.Source.UserID, message.ID, id, cw.ContentType, cw.Content)

			// Audio message
			case *linebot.AudioMessage:
				// TODO: AudioMessage

				cw, err := LineBot.GetMessageContent(message.ID).Do()
				if err != nil {
					log.Error("Get line file content", err)
					return
				}

				DiscordSendFile(event.Source.UserID, message.ID, id, cw.ContentType, cw.Content)

			// File message is not supported
			case *linebot.FileMessage:

				cw, err := LineBot.GetMessageContent(message.ID).Do()
				if err != nil {
					log.Error("Get line file content", err)
					return
				}

				DiscordSendFile(event.Source.UserID, message.ID, id, cw.ContentType, cw.Content)
			}

		}

	}
}

func getDiscordID(event *linebot.Event) string {
	var lid string
	var title string

	switch event.Source.Type {
	case linebot.EventSourceTypeUser:
		lid = event.Source.UserID
		profile, err := LineBot.GetProfile(lid).Do()
		if err != nil {
			log.Error("Get line profile", err)
		}
		title = "User | " + profile.DisplayName

	case linebot.EventSourceTypeGroup:
		lid = event.Source.GroupID
		profile, err := LineBot.GetGroupSummary(lid).Do()
		if err != nil {
			log.Error("Get Group Summary", err)
		}
		title = "Group | " + profile.GroupName

	case linebot.EventSourceTypeRoom:
		lid = event.Source.RoomID
		title = "Talk | " + lid

	}

	c := Channel{
		LineID: lid,
		Title:  title,
	}
	c.ByLineID()
	if c.Title != title {
		c.update(title)
	}

	return c.DiscordID
}
