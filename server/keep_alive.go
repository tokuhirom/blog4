package server

import (
	"log"
	"net/http"
	"time"
)

// KeepAlive apprun が今 min-scale 0 なので､寝ないように自分で自分を起こし続ける
func KeepAlive(url string) {
	log.Printf("Starting keep-alive process for %s...", url)
	for {
		time.Sleep(10 * time.Second)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("failed to request %s: %v", url, err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("unexpected status code: %d", resp.StatusCode)
		}
		_ = resp.Body.Close()
	}
}
