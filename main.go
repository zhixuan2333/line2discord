package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/zhixuan2333/line2discord/db"
)

//go:generate go run github.com/prisma/prisma-client-go generate

var (
	LineBot    *linebot.Client
	DiscordBot *discordgo.Session
	ctx        context.Context
	client     *db.PrismaClient
	PORT       string
)

var (
	LinechannelSecret string
	LinechannelToken  string
	DiscordToken      string
	GuildID           string
	ParentID          string
)

const preview = "?width=486&height=487"

func init() {
	err := godotenv.Load()
	if err != nil {
		Warm("Not found .env file passed", err)
	}

	GuildID = os.Getenv("GUILD_ID")
	ParentID = os.Getenv("PARENT_ID")
	LinechannelSecret = os.Getenv("LINE_CHANNEL_SECRET")
	LinechannelToken = os.Getenv("LINE_CHANNEL_TOKEN")
	DiscordToken = os.Getenv("DISCORD_TOKEN")

	PORT = ":" + os.Getenv("PORT")
	if PORT == ":" {
		PORT = ":8080"
	}

	if GuildID == "" ||
		ParentID == "" ||
		LinechannelSecret == "" ||
		LinechannelToken == "" ||
		DiscordToken == "" {
		Error("Not found env. \n(ex. GuildID, ParentID, LinechannelSecret, LinechannelToken, DiscordToken\n", nil)
		os.Exit(1)
	}

}

func main() {

	// Init Database
	ctx = context.Background()
	client = db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()
	Success("Database connected")

	// Init Line bot
	var err error
	LineBot, err = linebot.New(LinechannelSecret, LinechannelToken)
	if err != nil {
		panic(err)
	}
	Success("Line bot Online")

	// Init Discord bot
	DiscordBot, err = discordgo.New("Bot " + DiscordToken)
	if err != nil {
		panic(err)
	}
	DiscordBot.AddHandler(messageCreate)
	DiscordBot.Open()
	Success("Discord bot Online")

	// Init Web server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	http.HandleFunc("/webhook", WebHook)
	Success(fmt.Sprintf("Server Open at http://localhost%v", PORT))

	http.ListenAndServe(PORT, nil)

	// Wait here until CTRL-C or other term signal is received.
	Success("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc
	DiscordBot.Close()

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
		Error("Get Discord info", err)
		return
	}

	if st.ParentID != ParentID {
		return
	}

	lid := getLineID(m.ChannelID)
	if lid == "" {
		return
	}

	if m.Attachments != nil {
		for _, v := range m.Attachments {
			if v.Width != 0 && v.Height != 0 {
				ct := strings.Split(v.URL, ".")

				switch ct[len(ct)-1] {
				case "jpg", "jpeg", "png", "gif":
					_, err := LineBot.PushMessage(lid, linebot.NewImageMessage(v.URL, v.URL+preview)).Do()
					if err != nil {
						Error("Send line Image", err)
					}
					ToLine(lid, m.ChannelID, "image")

				case "mp4", "webm", "mkv", "flv", "avi", "mov", "wmv", "mpg", "mpeg":
					_, err := LineBot.PushMessage(lid, linebot.NewVideoMessage(v.URL, v.URL+preview)).Do()
					if err != nil {
						Error("Send line video", err)
					}
					ToLine(lid, m.ChannelID, "video")

				// if is not image or video then send url
				default:
					_, err := LineBot.PushMessage(lid, linebot.NewTextMessage(v.URL)).Do()
					if err != nil {
						Error("Send line file", err)
					}
					ToLine(lid, m.ChannelID, "file")
				}

			}
		}
	}
	if m.Content != "" {
		_, err = LineBot.PushMessage(lid, linebot.NewTextMessage(m.Content)).Do()
		if err != nil {
			Error("Send line message", err)
		}
		ToLine(lid, m.ChannelID, "message")

	}

}

// WebHook
func WebHook(w http.ResponseWriter, req *http.Request) {
	events, err := LineBot.ParseRequest(req)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {

		switch event.Type {
		case linebot.EventTypeMessage:

			switch message := event.Message.(type) {

			// Text message
			case *linebot.TextMessage:
				id := getDiscordID(event)

				DiscordSendMessage(event.Source.UserID, id, message.Text)

			// Image message
			case *linebot.ImageMessage:
				id := getDiscordID(event)

				cw, err := LineBot.GetMessageContent(message.ID).Do()
				if err != nil {
					Error("Get line file content", err)
				}

				// TODO: Auto get file extension
				ct := make([]string, 2)
				ct = strings.Split(cw.ContentType, "/")

				DiscordSendFile(event.Source.UserID, id, message.ID+"."+ct[1], cw.Content)

			// Video message
			case *linebot.VideoMessage:
				id := getDiscordID(event)

				cw, err := LineBot.GetMessageContent(message.ID).Do()
				if err != nil {
					Error("Get line file content", err)
				}

				// TODO: Auto get file extension
				ct := make([]string, 2)
				ct = strings.Split(cw.ContentType, "/")

				DiscordSendFile(event.Source.UserID, id, message.ID+"."+ct[1], cw.Content)

			// Audio message
			case *linebot.AudioMessage:
				// TODO: AudioMessage

				id := getDiscordID(event)

				cw, err := LineBot.GetMessageContent(message.ID).Do()
				if err != nil {
					Error("Get line file content", err)
				}

				// TODO: Auto get file extension
				ct := make([]string, 2)
				ct = strings.Split(cw.ContentType, "/")

				DiscordSendFile(event.Source.UserID, id, message.ID+"."+ct[1], cw.Content)

			// File message is not supported
			case *linebot.FileMessage:
				id := getDiscordID(event)

				cw, err := LineBot.GetMessageContent(message.ID).Do()
				if err != nil {
					Error("Get line file content", err)
				}

				DiscordSendFile(event.Source.UserID, id, message.ID+cw.ContentType, cw.Content)
			}

		}

	}
}
