package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"atomicmaya.me/misskey-ntfy-bridge/v2/models"
	"github.com/mitchellh/mapstructure"
)

func HandleNotePosted(apEvent map[string]any) {
	var body models.ActivityPubNoteEvent
	mapstructure.Decode(apEvent, &body)

	var toot string
	if len(body.Note.Text) >= 50 {
		toot = body.Note.Text[:50]
	} else {
		toot = body.Note.Text
	}

	if body.Note.ContentWarning != "" {
		toot = body.Note.ContentWarning
	}

	text := fmt.Sprintf(`%s (%s@%s) posted:
	%s...`,
		body.Note.User.Name,
		body.Note.User.UserName,
		body.Note.User.Host,
		toot,
	)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s", os.Getenv("NTFY_URL"), os.Getenv("NTFY_CHANNEL")), bytes.NewBufferString(text))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("NTFY_TOKEN")))
	req.Header.Set("Click", fmt.Sprintf("%s/notes/%s", os.Getenv("ORIGIN_URL"), body.Note.ID))

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()
}
