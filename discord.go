package main

import (
	"fmt"
	"io"

	"github.com/gabriel-vasile/mimetype"
	log "github.com/sirupsen/logrus"
)

func DiscordSendMessage(Author, channelID, message string) error {
	profile, err := LineBot.GetProfile(Author).Do()
	if err != nil {
		log.Error("Get line profile", err)
		return err
	}

	sm := fmt.Sprintf("%s: %s", profile.DisplayName, message)
	_, err = DiscordBot.ChannelMessageSend(channelID, sm)
	if err != nil {
		log.Error("Send message to discord", err)
	}
	ToDiscord(Author, channelID, "message")
	return nil
}

func DiscordSendFile(Author, messageID, channelID, ct string, cw io.ReadCloser) error {
	profile, err := LineBot.GetProfile(Author).Do()
	if err != nil {
		log.Error("Get line profile", err)
		return err
	}
	sm := fmt.Sprintf("%s:", profile.DisplayName)
	ext := mimetype.Lookup(ct)
	log.Info(ext.Extension())

	// TODO: Change to ChannelMessageSendComplex
	_, err = DiscordBot.ChannelFileSendWithMessage(channelID, sm, messageID+ext.Extension(), cw)
	if err != nil {
		log.Error("Send file to discord", err)
		return err
	}
	ToDiscord(Author, channelID, "file")

	return nil
}

func ToDiscord(lid, id, types string) {
	log.Infof("Send meesage to discord from: %v to: %v type: %v", lid, id, types)
}
