package main

import (
	"ads_redirector/lib/downloader"
	"ads_redirector/lib/reader"
	"ads_redirector/lib/round_robin"
	"github.com/buaazp/fasthttprouter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"gopkg.in/telegram-bot-api.v4"
)

const RESULT_FILE = "result.txt"

var r = round_robin.Stub()

func init() {
	viper.AddConfigPath("config")
	viper.ReadInConfig()
	log.SetFormatter(&log.TextFormatter{})
}

func Robin(ctx *fasthttp.RequestCtx) {
	next := r.Next()
	log.Infof("Redirected to [%s]", next)
	ctx.Redirect(next, 302)
}

func main() {
	go telegaInit()

	router := fasthttprouter.New()
	router.GET("/robin", Robin)

	log.Fatal(fasthttp.ListenAndServe(viper.GetString("APP_PORT"), router.Handler))
}

func telegaInit() {
	bot, err := tgbotapi.NewBotAPI(viper.GetString("TELEGRAM_KEY"))

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = viper.GetBool("TELEGRAM_DEBUG_MOD")

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if (update.Message.Document != nil) && (update.Message.Document.FileName == RESULT_FILE) {
			log.Printf("---------------------Get new result file----------------------------")
			downloader.GetFile(update.Message.Document.FileID, bot)
			changeRoundRobinUrl()
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ok, links updated")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}

func changeRoundRobinUrl() {
	links := reader.ReadFromFile(RESULT_FILE)
	r.New(links)
}
