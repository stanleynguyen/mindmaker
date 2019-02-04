package reducer

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (r *Reducer) handleListCommand(update tgbotapi.Update) {
	bucketName, err := r.Persistence.GetDefaultBucket(update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	} else if bucketName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "There's currently no bucket set to read from 😕")
		r.Bot.Send(msg)
		return
	}

	options, err := r.Persistence.ReadAllOptions(bucketName)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, getFormattedListOfOptions(options))
	r.Bot.Send(msg)
}