package model

import (
	"database/sql"
)

// User is the model for users.
type User struct {
	Id          int64          `json:"id"`
	Username    string         `json:"username"`
	Email       string         `json:"email"`
	Password    string         `json:"password,omitempty"`
	Location    sql.NullString `json:"location,omitempty"`
	IsValidated bool           `json:"isValidated"`
	CreatedAt   int64          `json:"createdAt"`
	ModifiedAt  int64          `json:"modifiedAt"`
}
