package sqlite

import "time"

func nowUTC() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func parseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
