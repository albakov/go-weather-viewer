package entity

import (
	"time"
)

type Session struct {
	Id        string
	UserId    int64
	ExpiresAt string
}

func (s *Session) ExpiresAtTime() time.Time {
	parsed, err := time.Parse(time.DateTime, s.ExpiresAt)
	if err != nil {
		return time.Time{}
	}

	return parsed
}
