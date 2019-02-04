package reducer

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stanleynguyen/mindmaker/domain"
)

func (r *Reducer) handleBucketsCommand(update tgbotapi.Update) {
	buckets, err := r.Persistence.GetAllBuckets(update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, getFormattedBucketsList(buckets))
	r.Bot.Send(msg)
}

func getFormattedBucketsList(buckets []domain.Bucket) string {
	if len(buckets) == 0 {
		return "Oops! You have not created any buckets yet."
	}

	rtv := "Here are the buckets that you've created 🤓:\n"
	for i, bucket := range buckets {
		rtv += fmt.Sprintf("- %v", bucket.Name)
		if i < len(buckets)-1 {
			rtv += "\n"
		}
	}

	return rtv
}
