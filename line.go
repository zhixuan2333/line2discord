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
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	for _, event := range events {
		c := getDiscordID(event)

		switch event.Type {
		case linebot.EventTypeMessage:

			switch message := event.Message.(type) {

			// Text message
			case *linebot.TextMessage:
				c.DiscordSendMessage(event.Source.UserID, message.Text)

			// Image message
			case *linebot.ImageMessage:
				c.DiscordSendFile(event.Source.UserID, message.ID)

			// Video message
			case *linebot.VideoMessage:
				c.DiscordSendFile(event.Source.UserID, message.ID)

			// Audio message (not work)
			case *linebot.AudioMessage:
				c.DiscordSendFile(event.Source.UserID, message.ID)

			// File message is not supported
			case *linebot.FileMessage:
				c.DiscordSendFile(event.Source.UserID, message.ID)

			case *linebot.StickerMessage:
				c.DiscordSendSticker(event.Source.UserID, message.PackageID+message.StickerID)

			}
		}

	}
}

func getDiscordID(event *linebot.Event) Channel {
	var lid string
	var title string

	switch event.Source.Type {
	case linebot.EventSourceTypeUser:
		lid = event.Source.UserID
		profile, err := LineBot.GetProfile(lid).Do()
		if err != nil {
			log.Error("Get line profile", err)
		}
		title = "👤・" + profile.DisplayName

	case linebot.EventSourceTypeGroup:
		lid = event.Source.GroupID
		profile, err := LineBot.GetGroupSummary(lid).Do()
		if err != nil {
			log.Error("Get Group Summary", err)
		}
		title = "👥・" + profile.GroupName

	case linebot.EventSourceTypeRoom:
		lid = event.Source.RoomID
		title = "🗣️・ " + lid

	}

	c := Channel{
		LineID: lid,
		Title:  title,
	}
	c.byLineID()
	if c.Title != title {
		c.update(title)
		_, err := discord.ChannelEdit(c.DiscordID, title)
		if err != nil {
			log.Errorf("Can't edit channel name: %v", err)
		}
	}
	return c
}
