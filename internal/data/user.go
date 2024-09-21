package data

import "time"

type User struct {
	Id          int64     `json:"id"`
	Username    string    `json:"username"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	Location    *string   `json:"location,omitempty"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Created     time.Time `json:"created"`
	LastSeen    time.Time `json:"lastSeen"`
}
