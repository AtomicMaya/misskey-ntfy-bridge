package main

import (
	"fmt"
	"net/http"
	"os"

	"atomicmaya.me/misskey-ntfy-bridge/v2/handlers"
	"atomicmaya.me/misskey-ntfy-bridge/v2/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf("Environment failed to load because of: %w", err)
	}
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"ping": "Alive"})
}

func fediEvent(c *gin.Context) {
	secret := c.Request.Header.Get("X-Misskey-Hook-Secret")
	if secret != os.Getenv("NTFY_TOKEN") {
		c.JSON(http.StatusUnauthorized, map[string]string{})
	}

	var event models.ActivityPubEvent
	if err := c.BindJSON(&event); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{})
	}

	if event.Server != os.Getenv("ORIGIN_URL") {
		c.JSON(http.StatusUnauthorized, map[string]string{})
	}

	fmt.Printf("TYPE %v\n", event.Type)
	fmt.Printf("BODY %v\n", event.Body)

	switch event.Type {
	case "followed":
		handlers.HandleFollowed(event.Body)
	case "follow":
		handlers.HandleFollow(event.Body)
	case "note":
		handlers.HandleNotePosted(event.Body)
	case "reply":
		handlers.HandleNoteReplied(event.Body)
	case "renote":
		handlers.HandleNoteBoosted(event.Body)
	}

	c.JSON(http.StatusOK, map[string]string{})
}

func main() {
	router := gin.Default()

	router.GET("/health", health)
	router.POST("/fedi-event", fediEvent)

	router.Run("0.0.0.0:1337")
}
