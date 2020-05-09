package main

import (
	parrotlib "CocoTelegramParrotBot/parrotlib"
	"github.com/davecgh/go-spew/spew"
	gocron "github.com/jasonlvhit/gocron"
	viper "github.com/spf13/viper"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var TelegramApiToken string = ""
var debug bool = false
var parrot parrotlib.Parrot
var bot *tgbotapi.BotAPI

func main() {

	getConfig()
	go setupCronjobs()

	spew.Config.Indent = "\t"

	var err error
	bot, err = tgbotapi.NewBotAPI(TelegramApiToken)
	if err != nil {
		log.Panic(spew.Sprintf("FATAL: error when creating telegram bot : %v", err))
	}

	bot.Debug = debug

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// new parrot
	parrot = parrotlib.NewParrot(
		"CoCo", "CoCo est cool, cool est CoCo !",
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

		// memorize the active channels
		parrot.AddChat(string(update.Message.Chat.ID))

		// add all talking users to the game
		parrot.AddUser(update.Message.From.UserName)

		if update.Message.IsCommand() {
			go handleCommand(bot, update.Message)
			continue
		}

		go handleMessage(bot, update.Message)
	}
}

func getConfig() {

	viper.SetDefault("tg_token", nil)
	viper.SetDefault("events_period_mins", 2)

	// handle config as file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/parrotTelegramParrotBot")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			log.Printf(spew.Sprintf("Error reading config file %v : %v", "viper config path", err))
		}
	}

	// handle config as env
	viper.SetEnvPrefix("coco") // will be uppercased automatically
	viper.BindEnv("tg_token")
	viper.BindEnv("events_period_mins")
	TelegramApiToken = viper.GetString("tg_token")

	log.Println(spew.Sprintf("TgToken: %v", TelegramApiToken))

	if len(TelegramApiToken) <= 0 {
		log.Panic("Null telegram API token given !")
	}
}

func handleCommand(bot *tgbotapi.BotAPI, command *tgbotapi.Message) {
	// Create a new MessageConfig. We don't have text yet,
	// so we should leave it empty.
	msg := tgbotapi.NewMessage(command.Chat.ID, "")

	// Extract the command from the Message.
	switch command.Command() {
	case "debug":
		msg.Text = spew.Sprintf("Debug: %v", parrot.Dump())
	case "help":
		msg.Text = "parrot the Parrot has no help, it's a f***ing parrot dude !\nAnyway try /silence, /target or /psst"
	case "silence":
		msg.Text = "Sileeence !"
	case "psst":
		msg.Text = "Roohh flying away !"
		randomShoulderSwitch()
	case "target":
		msg.Text = spew.Sprintf("Oooo i goooo to %v shoulder", command.CommandArguments())
		targetedUserName := command.CommandArguments()
		parrot.AddUser(targetedUserName)
		parrot.SwitchShoulder(targetedUserName)
	case "whereareyou":
		msg.Text = spew.Sprintf("I'am on %v shoulder ! heyk !", parrot.GetCurrentShoulder())
	default:
		msg.Text = "I doooon't knoooow that command"
	}

	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
}

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)
	parrot.AddUser(message.From.UserName)

	if message.Text == "" {
		return
	}

	if parrot.WillRepeat() {
		log.Printf("Parrot will repeat this time ...")
		msg := tgbotapi.NewMessage(message.Chat.ID, spew.Sprintf(parrot.Repeat(message.Text), *message.From))
		bot.Send(msg)
	} else {
		log.Printf("Parrot will be quiet for now ...")
	}
}

func triggerRandomEvent() {
	events := []interface{}{randomShoulderSwitch, randomSentence, randomQuizz, PreferedSentence}
	eventId := int(rand.New(rand.NewSource((time.Now().UnixNano()))).Float64() * float64(len(events)))

	log.Println(spew.Sprintf("Triggered random event: %v", events[eventId]))
	events[eventId].(func())()
}

func PreferedSentence() {
	log.Println(spew.Sprintf("PreferedSentence: starting"))
	// TODO: need a chat lists
	// msg := tgbotapi.NewMessage(, parrot.SayPreferedSentance())
	// if _, err := bot.Send(); err != nil {
	// 	log.Panic(err)
	// }
}

func randomShoulderSwitch() {
	log.Println(spew.Sprintf("randomShoulderSwitch: starting"))
	ruser, ok := parrot.RandomUser()
	if ok {
		parrot.SwitchShoulder(ruser)
		Broadcast(spew.Sprintf("Switched to %v shoulder ! heykk !!", ruser))
	}
}

func randomSentence() {
	log.Println(spew.Sprintf("randomSentence: starting"))
}

func randomQuizz() {
	log.Println(spew.Sprintf("randomQuizz: starting"))
}

func Broadcast(text string) {
	for _, chatId := range parrot.GetChats() {
		realChatId, _ := strconv.ParseInt(chatId, 10, 64)
		msg := tgbotapi.NewMessage(realChatId, text)
		bot.Send(msg)
	}
}

func setupCronjobs() {
	gocron.Every(viper.GetUint64("events_period_mins")).Minutes().Do(triggerRandomEvent)
	// Start all the pending jobs
	<-gocron.Start()
}
