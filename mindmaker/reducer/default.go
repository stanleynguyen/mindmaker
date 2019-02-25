package reducer

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const defaultResponse string = "Sorry boss 🙇 I didn't get what you just said\nBut I have a dad joke for you if you reply `yes` to this message 😎"

func (r *Reducer) handleAllOthersUpdate(update tgbotapi.Update) {
	if update.Message.Chat.IsGroup() || update.Message.Chat.IsSuperGroup() || update.Message.Chat.IsChannel() {
		if !r.handleUpdateFromGroup(update) {
			return
		}
	}

	isYesToDadJokeQn := update.Message.ReplyToMessage != nil &&
		update.Message.ReplyToMessage.Text == defaultResponse &&
		strings.EqualFold(update.Message.Text, "YES")
	if isYesToDadJokeQn {
		joke, err := getRandomDadJoke()
		var msg tgbotapi.MessageConfig
		if err != nil {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "dad joke")
		} else {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, joke)
		}

		r.Bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, defaultResponse)
	r.Bot.Send(msg)
}

func getRandomDadJoke() (string, error) {
	dadJokeServiceErr := errors.New("Dont have one on top of my head right now 😶 Sorry boss")

	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com", nil)
	if err != nil {
		return "", dadJokeServiceErr
	}
	req.Header.Add("Accept", "text/plain")
	httpsCli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{},
		},
	}
	resp, err := httpsCli.Do(req)
	if err != nil {
		return "", dadJokeServiceErr
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", dadJokeServiceErr
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	body := string(bodyBytes)
	return body, nil
}

// handleUpdateFromGroup take care of message from group
// returning whether to continue with dad joke
// if mindmaker is in the new members, send greeting message
// ignore all group events otherwise
func (r *Reducer) handleUpdateFromGroup(update tgbotapi.Update) bool {
	hasNewMembers := update.Message.NewChatMembers != nil
	if hasNewMembers {
		newMembers := *update.Message.NewChatMembers
		if len(newMembers) > 0 {
			for _, member := range newMembers {
				if member.ID == r.Bot.Self.ID {
					r.handleStartCommand(update)
					return false
				}
			}
		}
	}

	hasMembersLeft := update.Message.LeftChatMember != nil
	hasNewChatTitle := update.Message.NewChatTitle != ""
	hasNewChatPic := update.Message.NewChatPhoto != nil
	hasPicRemoval := update.Message.DeleteChatPhoto
	hasGroupCreation := update.Message.GroupChatCreated
	hasSuperGroupCreation := update.Message.SuperGroupChatCreated
	hasChannelCreation := update.Message.ChannelChatCreated

	continueWithDefaultFlow := !(hasNewMembers ||
		hasMembersLeft ||
		hasNewChatTitle ||
		hasNewChatPic ||
		hasPicRemoval ||
		hasGroupCreation ||
		hasSuperGroupCreation ||
		hasChannelCreation)

	return continueWithDefaultFlow
}
