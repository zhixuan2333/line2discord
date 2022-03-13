package main

import (
	"strings"
	"time"

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
	_, err = discord.ChannelMessageSendComplex(c.DiscordID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color: 0x5A65F1,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    profile.DisplayName,
				IconURL: profile.PictureURL,
				URL:     "https://example.com/#" + Author,
			},
			Description: message,
			Footer: &discordgo.MessageEmbedFooter{
				IconURL: discord.State.User.AvatarURL(""),
				Text:    discord.State.User.Username,
			},
			Timestamp: time.Now().Format(time.RFC3339),
		},
	})
	if err != nil {
		log.Error("Send message to discord", err)
		return err
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

	filetype := strings.Split(cw.ContentType, "/")[0]
	ext := mimetype.Lookup(cw.ContentType)

	message := &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   messageID + ext.Extension(),
				Reader: cw.Content,
			},
		},
		Embed: &discordgo.MessageEmbed{
			Color: 0x5A65F1,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    profile.DisplayName,
				IconURL: profile.PictureURL,
				URL:     "https://example.com/#" + Author,
			},
			Footer: &discordgo.MessageEmbedFooter{
				IconURL: discord.State.User.AvatarURL(""),
				Text:    discord.State.User.Username,
			},
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}

	switch filetype {
	case "image":
		message.Embed.Image = &discordgo.MessageEmbedImage{
			URL: "attachment://" + messageID + ext.Extension(),
		}
		message.Embed.Description = "Image send from LINE"

	case "video":
		message.Embed.Video = &discordgo.MessageEmbedVideo{
			URL: "attachment://" + messageID + ext.Extension(),
		}
		message.Embed.Description = "Video send from LINE"

	default:
		message.Embed.Description = "Some file send from Line."
	}

	_, err = discord.ChannelMessageSendComplex(c.DiscordID, message)
	if err != nil {
		log.Errorf("Send message to discord %v\n", err)
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

	st, err := discord.Channel(m.ChannelID)
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
