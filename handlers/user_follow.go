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

func HandleFollowEvent(apEvent map[string]any, outbound bool) {
	var body models.ActivityPubFollowEvent
	mapstructure.Decode(apEvent, &body)

	description := utils.SanitizeString(utils.Substring(body.User.Description, 250))

	var text string
	if !outbound {
		text = "New Follower: "
	} else {
		text = "Now Following: "
	}

	text += fmt.Sprintf(`%s (%s@%s)
%s...`,
		body.User.Name, body.User.UserName, body.User.Host,
		description,
	)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/%s", os.Getenv("NTFY_URL"), os.Getenv("NTFY_CHANNEL")), bytes.NewBufferString(utils.SanitizeString(text)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("NTFY_TOKEN")))
	req.Header.Set("Click", fmt.Sprintf("%s/@%s@%s", os.Getenv("ORIGIN_URL"), body.User.UserName, body.User.Host))

	client := &http.Client{}
	res, _ := client.Do(req)
	defer res.Body.Close()
}
