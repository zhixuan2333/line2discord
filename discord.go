package main

import (
	"fmt"
	"io"
)

func DiscordSendMessage(Author, channelID, message string) error {
	profile, err := LineBot.GetProfile(Author).Do()
	if err != nil {
		Error("Get line profile", err)
		return err
	}

	sm := fmt.Sprintf("%s: %s", profile.DisplayName, message)
	_, err = DiscordBot.ChannelMessageSend(channelID, sm)
	if err != nil {
		Error("Send message to discord", err)
	}
	ToDiscord(Author, channelID, "message")
	return nil
}

func DiscordSendFile(Author, channelID, filename string, file io.Reader) error {
	profile, err := LineBot.GetProfile(Author).Do()
	if err != nil {
		Error("Get line profile", err)
		return err
	}

	sm := fmt.Sprintf("%s:", profile.DisplayName)
	// TODO: Change to ChannelMessageSendComplex
	_, err = DiscordBot.ChannelFileSendWithMessage(channelID, sm, filename, file)
	if err != nil {
		Error("Send file to discord", err)
		return err
	}
	ToDiscord(Author, channelID, "file")

	return nil
}
