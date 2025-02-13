package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"atomicmaya.me/misskey-ntfy-bridge/v2/models"
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
	if len(body.Note.Text) >= 250 {
		toot = body.Note.Text[:250]
	} else {
		toot = body.Note.Text
	}

	if body.Note.ContentWarning != "" {
		toot = body.Note.ContentWarning
	}

	var text string

	switch eventType {
	case AP_POST:
		text = fmt.Sprintf(`%s (%s@%s) posted:
%s...`,
			body.Note.User.Name,
			body.Note.User.UserName,
			body.Note.User.Host,
			toot,
		)
	case AP_REPLY:
		var reply string
		if len(body.Note.Reply.Text) >= 250 {
			reply = body.Note.Reply.Text[:250]
		} else {
			reply = body.Note.Reply.Text
		}

		if body.Note.ContentWarning != "" {
			toot = body.Note.ContentWarning
			reply = body.Note.Reply.ContentWarning
		}

		text = fmt.Sprintf(`%s (%s@%s) replied:
> %s...
%s...`,
			body.Note.User.Name,
			body.Note.User.UserName,
			body.Note.User.Host,
			toot,
			reply,
		)
	case AP_BOOST:
		var boost string
		if len(body.Note.Text) >= 250 {
			boost = body.Note.Text[:250]
		} else {
			boost = body.Note.Text
		}

		var toot string
		if len(body.Note.Renote.Text) >= 250 {
			toot = body.Note.Renote.Text[:250]
		} else {
			toot = body.Note.Renote.Text
		}

		if body.Note.ContentWarning != "" {
			boost = body.Note.ContentWarning
			toot = body.Note.Renote.ContentWarning
		}

		text = fmt.Sprintf(`%s (%s@%s) boosted:
> %s...
%s...`,
			body.Note.User.Name,
			body.Note.User.UserName,
			body.Note.User.Host,
			toot,
			boost,
		)
	case AP_MENTION:
		text = fmt.Sprintf(`%s (%s@%s) mentioned you:
%s...`,
			body.Note.User.Name,
			body.Note.User.UserName,
			body.Note.User.Host,
			toot,
		)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s", os.Getenv("NTFY_URL"), os.Getenv("NTFY_CHANNEL")), bytes.NewBufferString(text))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("NTFY_TOKEN")))
	if eventType == AP_REPLY {
		req.Header.Set("Click", fmt.Sprintf("%s/notes/%s", os.Getenv("ORIGIN_URL"), body.Note.Reply.ID))
	} else if eventType == AP_BOOST {
		req.Header.Set("Click", fmt.Sprintf("%s/notes/%s", os.Getenv("ORIGIN_URL"), body.Note.Renote.ID))
	} else {
		req.Header.Set("Click", fmt.Sprintf("%s/notes/%s", os.Getenv("ORIGIN_URL"), body.Note.ID))
	}

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()
}
