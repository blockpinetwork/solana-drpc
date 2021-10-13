package utils

import "time"

func CurrentTimestamp() int64 {
	return  time.Now().Unix()
}

func CurrentTimestampMilli() int64 {
	return  time.Now().UnixNano() / int64(time.Millisecond)
}

func CurrentHourTimestamp() int64 {
	now := time.Now()
	timestamp := now.Unix() - int64(now.Second()) - int64(60 * now.Minute())
	return timestamp
}