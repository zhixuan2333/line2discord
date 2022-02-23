package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gabriel-vasile/mimetype"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	log "github.com/sirupsen/logrus"
)

func (c *Channel) DiscordSendMessage(Author, message string) error {
	profile, err := LineBot.GetProfile(Author).Do()
	if err != nil {
		log.Error("Get line profile", err)
		return err
	}

	sm := fmt.Sprintf("%s: %s", profile.DisplayName, message)
	_, err = DiscordBot.ChannelMessageSend(c.DiscordID, sm)
	if err != nil {
		log.Error("Send message to discord", err)
	}
	ToDiscord(Author, c.DiscordID, "message")
	return nil
}

func (c *Channel) DiscordSendFile(Author, messageID string) error {

	cw, err := LineBot.GetMessageContent(messageID).Do()
	if err != nil {
		log.Error("Get line file content: ", err)
		return err
	}

	profile, err := LineBot.GetProfile(Author).Do()
	if err != nil {
		log.Error("Get line profile", err)
		return err
	}
	sm := fmt.Sprintf("%s:", profile.DisplayName)
	ext := mimetype.Lookup(cw.ContentType)
	log.Info(ext.Extension())

	// TODO: Change to ChannelMessageSendComplex
	_, err = DiscordBot.ChannelFileSendWithMessage(c.DiscordID, sm, messageID+ext.Extension(), cw.Content)
	if err != nil {
		log.Error("Send file to discord", err)
		return err
	}
	ToDiscord(Author, c.DiscordID, "file")

	return nil
}

func ToDiscord(lid, id, types string) {
	log.Infof("Send meesage to discord from: %v to: %v type: %v", lid, id, types)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!") {

		return
	}

	st, err := DiscordBot.Channel(m.ChannelID)
	if err != nil {
		log.Error("Get Discord info", err)
		return
	}

	if st.ParentID != ParentID {
		return
	}

	var c Channel
	c.ByDiscordID(m.ChannelID)
	if c.LineID == "" {
		return
	}
	if m.Attachments != nil {
		for _, v := range m.Attachments {
			if v.Width != 0 && v.Height != 0 {
				ct := strings.Split(v.URL, ".")

				switch ct[len(ct)-1] {
				case "jpg", "jpeg", "png", "gif":
					_, err := LineBot.PushMessage(c.LineID, linebot.NewImageMessage(v.URL, v.URL+preview)).Do()
					if err != nil {
						log.Error("Send line Image", err)
					}
					ToLine(c.LineID, c.DiscordID, "image")

				case "mp4", "webm", "mkv", "flv", "avi", "mov", "wmv", "mpg", "mpeg":
					_, err := LineBot.PushMessage(c.LineID, linebot.NewVideoMessage(v.URL, v.URL+preview)).Do()
					if err != nil {
						log.Error("Send line video", err)
					}
					ToLine(c.LineID, c.DiscordID, "video")

				// if is not image or video then send url
				default:
					_, err := LineBot.PushMessage(c.LineID, linebot.NewTextMessage(v.URL)).Do()
					if err != nil {
						log.Error("Send line file", err)
					}
					ToLine(c.LineID, c.DiscordID, "file")
				}

			}
		}
	}
	if m.Content != "" {
		_, err = LineBot.PushMessage(c.LineID, linebot.NewTextMessage(m.Content)).Do()
		if err != nil {
			log.Error("Send line message", err)
		}
		ToLine(c.LineID, c.DiscordID, "message")

	}

}

func ToLine(lid, id, types string) {
	log.Infof("Send meesage to line from: %v to: %v type: %v", lid, id, types)
}
