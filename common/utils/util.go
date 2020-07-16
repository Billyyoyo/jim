package utils

import "time"

func GetCurrentMS() int64 {
	return time.Now().UnixNano() / 1e6
}
