package admin

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// NotifyHub sends PubSubHubbub notification to the Hub.
// hubURL: Hub URL (e.g.: https://pubsubhubbub.appspot.com/)
// feedURL: Feed URL (e.g.: https://example.com/feed)
func NotifyHub(hubURL string, feedURL string) error {
	// Create form values
	formData := url.Values{}
	formData.Set("hub.mode", "publish")
	formData.Set("hub.url", feedURL)

	// Create request
	resp, err := http.Post(
		hubURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		log.Printf("Failed to notify Hub: %v\n", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err := fmt.Errorf("hub notification failed: %d %s", resp.StatusCode, resp.Status)
		log.Printf("Failed to notify Hub: %v\n", err)
		return err
	}

	log.Printf("Notification sent to Hub: %s\n", hubURL)
	return nil
}
