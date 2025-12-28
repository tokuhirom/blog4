package utils

import "time"

func formatDateTime(time_ time.Time) string {
	return time_.Format("2006-01-02(Mon) 15:04")
}
