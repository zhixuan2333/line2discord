package main

import (
	"fmt"

	"github.com/gabriel-vasile/mimetype"
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
