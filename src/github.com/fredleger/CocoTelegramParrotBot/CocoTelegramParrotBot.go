package main

import (
	parrotlib "CocoTelegramParrotBot/parrotlib"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var TgToken string = ""
var debug bool = false

func getConfig() {
	// Telegram token
	TgToken = os.Getenv("TG_TOKEN")
	if len(TgToken) <= 0 {
		log.Panic("Invalid (null) telegram Token !")
	}
}

func main() {

	getConfig()

	spew.Config.Indent = "\t"

	bot, err := tgbotapi.NewBotAPI(TgToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// new parrot
	coco := parrotlib.NewParrot(
		"coco", "coco est cool, cool est coco !",
		"rrooohh ! %v a dit : ", 0.05, 17,
	)

	// waiting for messages
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// looping around all events
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			go handleCommand(&coco, bot, update.Message)
			continue
		}

		go handleMessage(&coco, bot, update.Message)
	}
}

func handleCommand(parrot *parrotlib.Parrot, bot *tgbotapi.BotAPI, command *tgbotapi.Message) {
	// Create a new MessageConfig. We don't have text yet,
	// so we should leave it empty.
	msg := tgbotapi.NewMessage(command.Chat.ID, "")

	// Extract the command from the Message.
	switch command.Command() {
	case "help":
		msg.Text = "Coco the Parrot has no help, it's a f***ing parrot dude !\nAnyway try /silence, /target or /psst"
	case "silence":
		msg.Text = "Sileeence !"
	case "psst":
		msg.Text = "Roohh flying away !"
	case "target":
		msg.Text = spew.Sprintf("Oooo i goooo to %v shoulder", command.CommandArguments())
		parrot.AddUser(command.CommandArguments())
		parrot.SwitchShoulder(command.CommandArguments())
	default:
		msg.Text = "I doooon't knoooow that command"
	}

	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func handleMessage(coco *parrotlib.Parrot, bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	coco.AddUser(message.From.UserName)

	if message.Text == "" {
		return
	}

	if coco.WillRepeat() {
		log.Printf("Parrot will repeat this time ...")
		msg := tgbotapi.NewMessage(message.Chat.ID, spew.Sprintf(coco.Repeat(message.Text), *message.From))
		bot.Send(msg)
	} else {
		log.Printf("Parrot will be quiet for now ...")
	}
}
