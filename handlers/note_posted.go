package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"atomicmaya.me/misskey-ntfy-bridge/v2/models"
	"atomicmaya.me/misskey-ntfy-bridge/v2/utils"
	"github.com/mitchellh/mapstructure"
)

type APEventType int

const (
	AP_POST    APEventType = 0
	AP_REPLY   APEventType = 1
	AP_BOOST   APEventType = 2
	AP_MENTION APEventType = 3
)

func HandleNotePosted(apEvent map[string]any, eventType APEventType) {
	var body models.ActivityPubNoteEvent
	mapstructure.Decode(apEvent, &body)

	var toot string
	if body.Note.ContentWarning != "" {
		toot = utils.Substring(body.Note.ContentWarning, 250)
	} else {
		toot = utils.Substring(body.Note.Text, 250)
	}

	var text string

	switch eventType {
	case AP_POST:
		text = fmt.Sprintf(`%s (%s@%s) posted:
%s...`,
			body.Note.User.Name, body.Note.User.UserName, body.Note.User.Host,
			toot,
		)
	case AP_REPLY:
		var reply string

		if body.Note.Reply.ContentWarning != "" {
			toot = utils.Substring(body.Note.ContentWarning, 250)
			reply = utils.Substring(body.Note.Reply.ContentWarning, 250)
		} else {
			reply = utils.Substring(body.Note.Reply.Text, 250)
		}

		text = fmt.Sprintf(`%s (%s@%s) replied:
> %s...
%s...`,
			body.Note.User.Name, body.Note.User.UserName, body.Note.User.Host,
			toot,
			reply,
		)
	case AP_BOOST:
		var boost string
		boost = toot[:]

		if body.Note.Renote.ContentWarning != "" {
			toot = utils.Substring(body.Note.Renote.ContentWarning, 250)
		} else {
			toot = utils.Substring(body.Note.Renote.Text, 250)
		}

		if len(boost) == 0 {
			text = fmt.Sprintf(`%s (%s@%s) boosted:
> %s...`,
				body.Note.User.Name, body.Note.User.UserName, body.Note.User.Host,
				toot,
			)
		} else {
			text = fmt.Sprintf(`%s (%s@%s) boosted:
> %s...
%s...`,
				body.Note.User.Name, body.Note.User.UserName, body.Note.User.Host,
				toot,
				boost,
			)
		}
	case AP_MENTION:
		text = fmt.Sprintf(`%s (%s@%s) mentioned you:
%s...`,
			body.Note.User.Name, body.Note.User.UserName, body.Note.User.Host,
			toot,
		)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s", os.Getenv("NTFY_URL"), os.Getenv("NTFY_CHANNEL")), bytes.NewBufferString(utils.SanitizeString(text)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("NTFY_TOKEN")))
	if eventType == AP_REPLY {
		req.Header.Set("Click", fmt.Sprintf("%s/notes/%s", os.Getenv("ORIGIN_URL"), body.Note.Reply.ID))
	} else {
		req.Header.Set("Click", fmt.Sprintf("%s/notes/%s", os.Getenv("ORIGIN_URL"), body.Note.ID))
	}

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()
}
