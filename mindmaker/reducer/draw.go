package reducer

import (
	"fmt"
	"log"
	"math/rand"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (r *Reducer) handleDrawCommand(update tgbotapi.Update) {
	bucketName, err := r.Persistence.GetDefaultBucket(update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	} else if bucketName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "There's currently no bucket set to draw a decision from 😕")
		r.Bot.Send(msg)
		return
	}

	options, err := r.Persistence.ReadAllOptions(bucketName)
	drawnOption := options[rand.Intn(len(options))]
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("🎊 Boss you have drawn decision %v 🎉", drawnOption))
	r.Bot.Send(msg)
}

func (r *Reducer) sendErrMessage(chatID int64) {
	errMsg := tgbotapi.NewMessage(chatID, "Sorry I'm not feeeling very well :( Please try again later")
	r.Bot.Send(errMsg)
}
