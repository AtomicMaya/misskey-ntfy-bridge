package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"atomicmaya.me/misskey-ntfy-bridge/v2/models"
	"github.com/mitchellh/mapstructure"
)

func HandleNoteBoosted(apEvent map[string]any) {
	var body models.ActivityPubNoteEvent
	mapstructure.Decode(apEvent, &body)

	var boost string
	if len(body.Note.Text) >= 50 {
		boost = body.Note.Text[:50]
	} else {
		boost = body.Note.Text
	}

	var toot string
	if len(body.Note.Renote.Text) >= 50 {
		toot = body.Note.Renote.Text[:50]
	} else {
		toot = body.Note.Renote.Text
	}

	if body.Note.ContentWarning != "" {
		boost = body.Note.ContentWarning
		toot = body.Note.Renote.ContentWarning
	}

	text := fmt.Sprintf(`%s (%s@%s) boosted:
	> %s...
	%s...`,
		body.Note.User.Name,
		body.Note.User.UserName,
		body.Note.User.Host,
		toot,
		boost,
	)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s", os.Getenv("NTFY_URL"), os.Getenv("NTFY_CHANNEL")), bytes.NewBufferString(text))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("NTFY_TOKEN")))
	req.Header.Set("Click", fmt.Sprintf("%s/notes/%s", os.Getenv("ORIGIN_URL"), body.Note.ID))

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()
}
