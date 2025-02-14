package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"

	"dev.catgirl.global/misskey-ntfy-bridge/v2/app/handlers"
	"dev.catgirl.global/misskey-ntfy-bridge/v2/app/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const (
	VERSION = "1.0.2"
	AUTHOR  = "maya <dev@catgirl.global>"
	SOURCE  = "https://github.com/AtomicMaya/misskey-ntfy-bridge"
	LICENSE = "EUPL-1.2"
)

// Inferred from ORIGIN_URL
var ORIGIN_DOMAIN string

// Loads environment variables if supplied via a .env configuration file.
func init() {
	// If the SOURCE flag is provided (e.g. via the Dockerfile build) then don't throw an error.
	if os.Getenv("SOURCE") == "container" {
		return
	}

	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", fmt.Errorf("environment failed to load because of: %w", err))
	}
}

// Checks for any missing environment keys which should be supplied to the program by means of either .env file, systemd configuration or some other docker file.
func init() {
	env_keys := []string{"HOST", "PORT", "ORIGIN_URL", "NTFY_URL", "NTFY_CHANNEL", "NTFY_TOKEN"}
	env_missing := []string{}
	for _, key := range env_keys {
		if os.Getenv(key) == "" {
			env_missing = append(env_missing, key)
		}
	}

	if len(env_missing) > 0 {
		fmt.Fprintf(os.Stderr, "error: %v\n", fmt.Errorf("environment variables missing: %v", env_missing))
		os.Exit(1)
	}
}

func init() {
	r := regexp.MustCompile(`^https?:\/\/`)
	ORIGIN_DOMAIN = r.ReplaceAllString(os.Getenv("ORIGIN_URL"), "")
}

// Endpoint function to be used by uptime monitoring software, e.g. uptime-kuma
func health(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "alive", "version": VERSION, "author": AUTHOR, "source": SOURCE, "license": LICENSE})
}

// Endpoint function triggered by the Misskey/Sharkey webhook
func fediEvent(c *gin.Context) {
	// Additional cross-instance malicious pollution protection
	if c.Request.Host != ORIGIN_DOMAIN {
		c.JSON(http.StatusUnauthorized, map[string]string{})
	}

	// Secret defined in the webhook interface, pass on to ntfy
	secret := c.Request.Header.Get("X-Misskey-Hook-Secret")

	// Limit the ability for malicious users to spam the webhook
	if secret != os.Getenv("NTFY_TOKEN") {
		c.JSON(http.StatusUnauthorized, map[string]string{})
		return
	}

	// Deserialize the top-level object
	var event models.ActivityPubEvent
	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{})
		return
	}

	// Additional cross-instance malicious pollution protection
	if event.Server != ORIGIN_DOMAIN {
		c.JSON(http.StatusUnauthorized, map[string]string{})
		return
	}

	var err error

	/*  Subtypes of ActivityPub events sent by Misskey/Sharkey
	As the data structure passed on is very similar in some cases,
	to avoid repetition we make use of common handlers with additional parameters
	*/
	switch event.Type {
	case "followed":
		err = handlers.HandleFollowEvent(event.Body, false)
	case "follow":
		err = handlers.HandleFollowEvent(event.Body, true)
	case "note":
		err = handlers.HandleNotePosted(event.Body, handlers.AP_POST)
	case "reply":
		err = handlers.HandleNotePosted(event.Body, handlers.AP_REPLY)
	case "renote":
		err = handlers.HandleNotePosted(event.Body, handlers.AP_BOOST)
	case "mention":
		err = handlers.HandleNotePosted(event.Body, handlers.AP_MENTION)
	case "reaction":
		// TODO: Not yet implemented on Sharkey
		fallthrough
	default:
		// If it is anything else, we don't know what to do with it and should simply reject it.
		c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": fmt.Sprintf("unsupported activitypub event type: '%s'", event.Type)})
		return
	}

	if err != nil {
		// Avoid leaking process information to potential attackers.
		c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("an unexpected error occured. Please contact the system administrator. If you are the system administrator, check the processes' logs ;)")})
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	} else {
		// Return an empty response, but 201 indicating that it has been accepted
		c.JSON(http.StatusAccepted, map[string]string{})
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/health", health)
	router.POST("/fedi-event", fediEvent)

	router.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
