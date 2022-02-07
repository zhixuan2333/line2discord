package main

import (
	"github.com/line/line-bot-sdk-go/v7/linebot"
	log "github.com/sirupsen/logrus"
)

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
