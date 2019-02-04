package reducer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/stanleynguyen/mindmaker/domain"
)

// BucketNameSeparator string used for concatting chatID with user defined name
// to form bucket names
const BucketNameSeparator = " - "

func getPrettyArgumentString(rawArgString string) string {
	return strings.Trim(rawArgString, " ")
}

func getBucketNameFromChatID(chatID int64, userGivenName string) string {
	return strconv.Itoa(int(chatID)) + BucketNameSeparator + userGivenName
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
