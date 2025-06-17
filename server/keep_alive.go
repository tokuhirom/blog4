package server

import (
	"log/slog"
	"net/http"
	"time"
)

// KeepAlive apprun が今 min-scale 0 なので､寝ないように自分で自分を起こし続ける
func KeepAlive(url string) {
	slog.Info("Starting keep-alive process", slog.String("url", url))
	for {
		time.Sleep(10 * time.Second)
		resp, err := http.Get(url)
		if err != nil {
			slog.Error("failed to request keep-alive URL", slog.String("url", url), slog.Any("error", err))
			continue
		}
		if resp.StatusCode != http.StatusOK {
			slog.Warn("unexpected keep-alive status code", slog.Int("status", resp.StatusCode), slog.String("url", url))
		}
		_ = resp.Body.Close()
	}
}
