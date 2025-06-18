package admin

import (
	"fmt"
	"io"
	"log/slog"
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
		slog.Error("Failed to notify Hub", slog.String("hub_url", hubURL), slog.String("feed_url", feedURL), slog.Any("error", err))
		return fmt.Errorf("failed to post to hub %s: %w", hubURL, err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err := fmt.Errorf("hub notification failed: %d %s", resp.StatusCode, resp.Status)
		slog.Error("Failed to notify Hub", slog.String("hub_url", hubURL), slog.Int("status_code", resp.StatusCode), slog.String("status", resp.Status))
		return fmt.Errorf("hub notification failed: %w", err)
	}

	slog.Info("Notification sent to Hub", slog.String("hub_url", hubURL), slog.String("feed_url", feedURL))
	return nil
}
