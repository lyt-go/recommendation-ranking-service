package model

import (
	"strings"
	"time"
)

type UserProfile struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Interests []string  `json:"interests"`
	Tags      []string  `json:"tags"`
	Region    string    `json:"region"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserProfile) Validate() error {
	u.UserID = strings.TrimSpace(u.UserID)
	u.Region = strings.TrimSpace(u.Region)
	if u.UserID == "" {
		return NewValidationError("user_id", "用户ID不能为空")
	}
	return nil
}

type UserProfileFilter struct {
	Region string
	Tag    string
}

func (f UserProfileFilter) Match(u *UserProfile) bool {
	if f.Region != "" && u.Region != f.Region {
		return false
	}
	if f.Tag != "" {
		found := false
		for _, t := range u.Tags {
			if t == f.Tag {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
