package reducer

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

func (r *Reducer) handleStartCommand(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		"Hey there~ Nice to meet you 😊 Let me tell you a story about the last time my friends had a dilemma of choice... Oops not any that I know of since I'm here 😝 As my name suggested, my purpose of existence is to help you with making up your mind 😉 So let's be friends and make decisions together (even though I do most of the work)",
	)
	r.Bot.Send(msg)
}
