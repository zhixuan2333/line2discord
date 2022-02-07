package main

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;not null;primary_key"`
	Title     string
	LineID    string `gorm:"type:varchar(100);not null"`
	DiscordID string `gorm:"type:varchar(100);not null"`
}

func (c *Channel) ByLineID() {
	result := db.Where(&Channel{LineID: c.LineID}).First(&c)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.createChannel()

		return
	} else if result.Error != nil {
		log.Error("get channel by line id", result.Error)

		return
	}
}

func (c *Channel) ByDiscordID(DiscordID string) {
	db.Where(&Channel{DiscordID: DiscordID}).First(&c)
}

func (c *Channel) update(title string) {
	db.Model(&c).Update("title", title)
}

func (c *Channel) createChannel() {
	channel, err := DiscordBot.GuildChannelCreateComplex(GuildID, discordgo.GuildChannelCreateData{
		Name:     c.Title,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: ParentID,
	})
	if err != nil {
		log.Error("create discord channel", err)
		return
	}

	c.ID = uuid.New()
	c.DiscordID = channel.ID

	db.Create(c)
}
