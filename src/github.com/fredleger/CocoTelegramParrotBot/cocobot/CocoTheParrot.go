package main

import (
    "os"
    "log"
    "gopkg.in/telegram-bot-api.v4"
    "github.com/fredleger/CocoTelegramParrotBot/parrotlib"
    // "github.com/davecgh/go-spew/spew"
)

func main() {

    var TgToken string
    TgToken = os.Getenv("TG_TOKEN")
    if (len(TgToken) <= 0) {
        log.Panic("Invalid (null) telegram Token !")
    }

    bot, err := tgbotapi.NewBotAPI(TgToken)
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = false

    log.Printf("Authorized on account %s", bot.Self.UserName)

    // new parrot
    coco := parrot.NewParrot("coco", "coco est cool, cool est coco !", "yeak !", 0.05)
    //coco.Dump()

    // waiting for messages
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }

        log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
        coco.AddUser(update.Message.From.UserName)

        if coco.WillRepeat() {
            log.Printf("Parrot is awake !!!")
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, coco.Repeat(update.Message.Text))
            bot.Send(msg)
        } else {
            log.Printf("Parrot will be quiet for now ...")
        }
    }
}


