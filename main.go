package main

import (
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	LineBot    *linebot.Client
	DiscordBot *discordgo.Session
	db         *gorm.DB
)

var (
	LinechannelSecret string
	LinechannelToken  string
	DiscordToken      string
	GuildID           string
	ParentID          string
	DatabaseURL       string
)

const preview = "?width=486&height=487"

func init() {
	log.SetLevel(log.InfoLevel)
	err := godotenv.Load()
	if err != nil {
		log.Warn("Not found .env file passed", err)
	}

	GuildID = os.Getenv("GUILD_ID")
	ParentID = os.Getenv("PARENT_ID")
	LinechannelSecret = os.Getenv("LINE_CHANNEL_SECRET")
	LinechannelToken = os.Getenv("LINE_CHANNEL_TOKEN")
	DiscordToken = os.Getenv("DISCORD_TOKEN")
	DatabaseURL = os.Getenv("DATABASE_URL")

	if GuildID == "" ||
		ParentID == "" ||
		LinechannelSecret == "" ||
		LinechannelToken == "" ||
		DiscordToken == "" ||
		DatabaseURL == "" {
		log.Panicf("Not found env. \n(ex. GuildID, ParentID, LinechannelSecret, LinechannelToken, DiscordToken\n", nil)
	}

}

func main() {
	var err error

	// Init Database
	// db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db, err = gorm.Open(postgres.Open(DatabaseURL), &gorm.Config{})

	if err != nil {
		log.Panic("Init Database", err)
	}
	db.AutoMigrate(&Channel{})
	log.Info("Successfully connencted to database")

	// Init Line bot
	LineBot, err = linebot.New(LinechannelSecret, LinechannelToken)
	if err != nil {
		panic(err)
	}
	log.Info("Successfully online line bot")

	// Init Discord bot
	DiscordBot, err = discordgo.New("Bot " + DiscordToken)
	if err != nil {
		panic(err)
	}
	DiscordBot.AddHandler(messageCreate)
	DiscordBot.Open()
	log.Info("Successfully online discord bot")

	PORT := ":" + os.Getenv("PORT")
	if PORT == ":" {
		PORT = ":8080"
	}
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
	log.Infof("Send meesage to line from: %v to: %v type: %v", lid, id, types)
}
