package reducer

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (r *Reducer) handleTakeoutCommand(update tgbotapi.Update) {
	defaultSet, err := r.Persistence.DefaultWasSet(update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}
	if !defaultSet {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "There's currently no bucket set to take out any decision from 😕")
		r.Bot.Send(msg)
		return
	}
	bucketName, err := r.Persistence.GetDefaultBucket(update.Message.Chat.ID)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}

	argStr := getPrettyArgumentString(update.Message.CommandArguments())
	options, err := r.Persistence.ReadAllOptions(bucketName)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}
	optionNumber, err := strconv.Atoi(argStr)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Boss 😳 can you please gimme a number of the option in this form: /takeout <decision number>?\nIf you forgot what the decision indexes are, here is a reminder:\n%v", getFormattedListOfOptions(options)))
		r.Bot.Send(msg)
		return
	} else if optionNumber > len(options) || optionNumber < 1 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Your option must be in range of 1-%v 🙃", len(options)))
		r.Bot.Send(msg)
		return
	}

	optionIdx := optionNumber - 1
	updatedOptions, err := r.Persistence.DeleteOption(bucketName, optionIdx)
	if err != nil {
		log.Println(err)
		r.sendErrMessage(update.Message.Chat.ID)
		return
	}
	userReadableBucketName := strings.SplitN(bucketName, BucketNameSeparator, 2)[1]
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Decision %v taken out of bucket %v 🤫 Here is the new list of decisions:\n%v", optionNumber, userReadableBucketName, getFormattedListOfOptions(updatedOptions)))
	r.Bot.Send(msg)
}
