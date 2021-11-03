package main

import (
	"errors"

	"github.com/zhixuan2333/line2discord/db"
)

func getRecordByLineID(lineID string) (*db.ChannelModel, error) {

	channel, err := client.Channel.FindUnique(
		db.Channel.LineID.Equals(lineID),
	).Exec(ctx)
	if errors.Is(err, db.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		Error("Get Record by lineID", err)
		return channel, err
	}

	return channel, nil

}

func getRecordByDiscordID(DiscordID string) (*db.ChannelModel, error) {

	channel, err := client.Channel.FindUnique(
		db.Channel.DiscordID.Equals(DiscordID),
	).Exec(ctx)
	if errors.Is(err, db.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		Error("Get Record by DiscordID", err)

		return channel, err
	}

	return channel, nil

}

func updateChannel(lineID, title string) (*db.ChannelModel, error) {
	channel, err := client.Channel.FindUnique(
		db.Channel.LineID.Equals(lineID),
	).Update(
		db.Channel.Title.Set(title),
	).Exec(ctx)
	if err != nil {
		Error("Update record title", err)

	}

	return channel, err
}

func createChannel(lineID, discordID, Title string) (*db.ChannelModel, error) {

	channel, err := client.Channel.CreateOne(
		db.Channel.Title.Set(Title),
		db.Channel.LineID.Set(lineID),
		db.Channel.DiscordID.Set(discordID),
	).Exec(ctx)
	if err != nil {
		Error("Create Record", err)
		return nil, err
	}
	return channel, nil
}
