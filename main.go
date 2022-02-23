package main

import (
	"net/http"
	"os"
	"os/signal"
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
