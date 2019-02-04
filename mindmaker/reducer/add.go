package reducer

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stanleynguyen/mindmaker/domain"
)

func (r *Reducer) handleAddCommand(update tgbotapi.Update) {
	bucketName, err := r.Persistence.GetDefaultBucket(update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	} else if bucketName == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "There's currently no bucket set to add this decision to 😕")
		r.Bot.Send(msg)
		return
	}

	argStr := getPrettyArgumentString(update.Message.CommandArguments())
	err = r.Persistence.InsertOption(bucketName, domain.Option(argStr))
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Decision %v added to bucket %v", argStr, bucketName))
	r.Bot.Send(msg)
}
