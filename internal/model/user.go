package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID
	Name  string
	Roles []string
}

func (u User) CheckRoles(roles []string) bool {
	role := make(map[string]struct{})
	for _, r := range roles {
		role[r] = struct{}{}
	}
	for _, ur := range u.Roles {
		if _, ok := role[ur]; ok {
			return true
		}
	}
	return false
}
