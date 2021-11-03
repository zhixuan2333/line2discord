package main

import (
	"github.com/line/line-bot-sdk-go/linebot"
)

func getDiscordID(event *linebot.Event) string {
	var lid string
	var title string

	switch event.Source.Type {
	case linebot.EventSourceTypeUser:
		lid = event.Source.UserID
		profile, err := LineBot.GetProfile(lid).Do()
		if err != nil {
			Error("Get line profile", err)
		}
		title = "User | " + profile.DisplayName

	case linebot.EventSourceTypeGroup:
		lid = event.Source.GroupID
		profile, err := LineBot.GetGroupSummary(lid).Do()
		if err != nil {
			Error("Get Group Summary", err)
		}
		title = "Group | " + profile.GroupName

	case linebot.EventSourceTypeRoom:
		lid = event.Source.RoomID
		title = "Talk | " + lid

	}

	channel, _ := getRecordByLineID(lid)
	if channel == nil {
		id, _ := DiscordcreateChannel(title)
		channel, _ = createChannel(lid, id, title)
	}

	if channel.Title != title {
		channel, _ = updateChannel(lid, title)
	}
	return channel.DiscordID
}

func getLineID(DiscordID string) string {
	channel, _ := getRecordByDiscordID(DiscordID)
	return channel.LineID

}
