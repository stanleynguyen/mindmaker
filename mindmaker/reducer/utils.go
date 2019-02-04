package reducer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/stanleynguyen/mindmaker/domain"
)

const BUCKET_NAME_SEPARATOR = " - "

func getPrettyArgumentString(rawArgString string) string {
	return strings.Trim(rawArgString, " ")
}

func getBucketNameFromChatID(chatID int64, userGivenName string) string {
	return strconv.Itoa(int(chatID)) + BUCKET_NAME_SEPARATOR + userGivenName
}

func getFormattedListOfOptions(options []domain.Option) string {
	if len(options) == 0 {
		return "Opps! Your current bucket doesnt contain any options"
	}

	rtv := "Possible decisions in your current default bucket:\n"
	for i := 0; i < len(options); i++ {
		rtv += fmt.Sprintf("%v. %v", i+1, options[i])
		if i < len(options)-1 {
			rtv += "\n"
		}
	}

	return rtv
}
