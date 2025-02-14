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

func HandleFollowEvent(apEvent map[string]any, outbound bool) error {
	// Deserialize the top level object
	var body models.ActivityPubFollowEvent

	if err := mapstructure.Decode(apEvent, &body); err != nil {
		return err
	}

	// Get the target user's description.
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

	// Generate the request to ntfy
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", os.Getenv("NTFY_URL"), os.Getenv("NTFY_CHANNEL")), bytes.NewBufferString(utils.SanitizeString(text)))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("NTFY_TOKEN")))
	req.Header.Set("Click", fmt.Sprintf("%s/@%s@%s", os.Getenv("ORIGIN_URL"), body.User.UserName, body.User.Host))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	return nil
}
