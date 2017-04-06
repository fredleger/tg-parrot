package main

import (
    "log"
    "gopkg.in/telegram-bot-api.v4"
    "github.com/fredleger/golang/parrot"
)

func main() {

    bot, err := tgbotapi.NewBotAPI("367497906:AAFz_E_iv11qVLboD-dSRhP8QK0ii8oVlro")
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = false

    log.Printf("Authorized on account %s", bot.Self.UserName)

    // new parrot
    coco := parrot.NewParrot("coco", "coco est cool, cool est coco !", "yeak !", 7)
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

        if coco.WillRepeat() {
            log.Printf("Parrot is awake !!!")
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, coco.Repeat(update.Message.Text))
            bot.Send(msg)
        } else {
            log.Printf("Parrot will be quiet for now ...")
        }
    }
}


