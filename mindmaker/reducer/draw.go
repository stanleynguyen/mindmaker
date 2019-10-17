package reducer

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (r *Reducer) handleDrawCommand(update tgbotapi.Update) {
	defaultSet, err := r.Persistence.DefaultWasSet(update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}
	if !defaultSet {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "There's currently no bucket set to draw a decision from 😕")
		r.Bot.Send(msg)
		return
	}
	bucketName, err := r.Persistence.GetDefaultBucket(update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	options, err := r.Persistence.ReadAllOptions(update.Message.Chat.ID, bucketName)
	rand.Seed(time.Now().UnixNano())
	drawnOption := options[rand.Intn(len(options))]
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("🎊 Boss you have drawn decision %v 🎉", drawnOption))
	r.Bot.Send(msg)
}
