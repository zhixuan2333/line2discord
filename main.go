package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//go:generate go run github.com/prisma/prisma-client-go generate

var (
	LineBot    *linebot.Client
	DiscordBot *discordgo.Session
	ctx        context.Context
	PORT       string
	db         *gorm.DB
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
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: false,
	})
	err := godotenv.Load()
	if err != nil {
		log.Warn("Not found .env file passed", err)
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
		log.Panicf("Not found env. \n(ex. GuildID, ParentID, LinechannelSecret, LinechannelToken, DiscordToken\n", nil)
	}

}

func main() {
	var err error

	// Init Database
	ctx = context.Background()
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Panic("Init Database", err)
	}
	db.AutoMigrate(&Channel{})

	// Init Line bot
	LineBot, err = linebot.New(LinechannelSecret, LinechannelToken)
	if err != nil {
		panic(err)
	}
	log.Info("Line bot Online")

	// Init Discord bot
	DiscordBot, err = discordgo.New("Bot " + DiscordToken)
	if err != nil {
		panic(err)
	}
	DiscordBot.AddHandler(messageCreate)
	DiscordBot.Open()
	log.Info("Discord bot Online")

	// Init Web server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	http.HandleFunc("/webhook", WebHook)
	log.Infof("Server Open at http://localhost%v", PORT)

	http.ListenAndServe(PORT, nil)

	// Wait here until CTRL-C or other term signal is received.
	log.Info("Bot is now running.  Press CTRL-C to exit.")
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
	log.Infof("[MESSAGE] | %33s | <-- | %18s | %7s |", lid, id, types)
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
					log.Error("Get line file content", err)
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
					log.Error("Get line file content", err)
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
					log.Error("Get line file content", err)
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
					log.Error("Get line file content", err)
				}

				DiscordSendFile(event.Source.UserID, id, message.ID+cw.ContentType, cw.Content)
			}

		}

	}
}
