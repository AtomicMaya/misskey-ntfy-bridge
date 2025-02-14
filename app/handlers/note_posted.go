package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"dev.catgirl.global/misskey-ntfy-bridge/v2/app/models"
	"dev.catgirl.global/misskey-ntfy-bridge/v2/app/utils"
	"github.com/mitchellh/mapstructure"
)

type APNoteEventType int

// Emum for APNoteEventType
const (
	AP_POST    APNoteEventType = 0
	AP_REPLY   APNoteEventType = 1
	AP_BOOST   APNoteEventType = 2
	AP_MENTION APNoteEventType = 3
)

// Handler for any form of note-related ActivityPub event
func HandleNotePosted(apEvent map[string]any, eventType APNoteEventType) error {
	// Deserialize the top level object
	var body models.ActivityPubNoteEvent
	if err := mapstructure.Decode(apEvent, &body); err != nil {
		return err
	}

	// Extract the toot contents, substituting with the content warning should one exist.
	// Limit the size to 250 characters to avoid overloading ntfy.
	var toot string
	if body.Note.ContentWarning != "" {
		toot = utils.Substring(body.Note.ContentWarning, 250)
	} else {
		toot = utils.Substring(body.Note.Text, 250)
	}

	var text string

	switch eventType {
	// Posts made by the user.
	case AP_POST:
		text = fmt.Sprintf(`%s (%s@%s) posted:
%s...`,
			body.Note.User.Name, body.Note.User.UserName, body.Note.User.Host,
			toot,
		)
	// Replies made to a user's post.
	// The reply may contain a content warning.
	case AP_REPLY:
		var reply string

		if body.Note.Reply.ContentWarning != "" {
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
	// Boosts made by other users of one of the user's posts.
	// Boosts may contain text, which may be cw'ed
	case AP_BOOST:
		var boost string
		// In this case, the boost is inverted from the note, so we copy the content over.
		boost = toot[:]

		if body.Note.Renote.ContentWarning != "" {
			toot = utils.Substring(body.Note.Renote.ContentWarning, 250)
		} else {
			toot = utils.Substring(body.Note.Renote.Text, 250)
		}

		text = fmt.Sprintf(`%s (%s@%s) boosted:
> %s...`,
			body.Note.User.Name, body.Note.User.UserName, body.Note.User.Host,
			toot,
		)

		// If the boost contains text.
		if len(boost) != 0 {
			text += fmt.Sprintf(`
%s...`,
				boost,
			)
		}
	// If the user has been mentioned.
	case AP_MENTION:
		text = fmt.Sprintf(`%s (%s@%s) mentioned you:
%s...`,
			body.Note.User.Name, body.Note.User.UserName, body.Note.User.Host,
			toot,
		)
	}

	// Generate the request to ntfy
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", os.Getenv("NTFY_URL"), os.Getenv("NTFY_CHANNEL")), bytes.NewBufferString(utils.SanitizeString(text)))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("NTFY_TOKEN")))
	if eventType == AP_REPLY {
		req.Header.Set("Click", fmt.Sprintf("%s/notes/%s", os.Getenv("ORIGIN_URL"), body.Note.Reply.ID))
	} else {
		req.Header.Set("Click", fmt.Sprintf("%s/notes/%s", os.Getenv("ORIGIN_URL"), body.Note.ID))
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return nil
}
