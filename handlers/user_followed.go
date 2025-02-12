package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"atomicmaya.me/misskey-ntfy-bridge/v2/models"
	"github.com/mitchellh/mapstructure"
)

func HandleFollowed(apEvent map[string]any) {
	var body models.ActivityPubFollowEvent
	mapstructure.Decode(apEvent, &body)

	var description string
	if len(body.User.Description) >= 50 {
		description = body.User.Description[:50]
	} else {
		description = body.User.Description
	}

	text := fmt.Sprintf(`New Follower: %s (%s@%s)
	%s...`,
		body.User.Name,
		body.User.UserName,
		body.User.Host,
		description,
	)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s", os.Getenv("NTFY_URL"), os.Getenv("NTFY_CHANNEL")), bytes.NewBufferString(text))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("NTFY_TOKEN")))
	req.Header.Set("Click", fmt.Sprintf("%s/@%s@%s", os.Getenv("ORIGIN_URL"), body.User.UserName, body.User.Host))

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()
}
