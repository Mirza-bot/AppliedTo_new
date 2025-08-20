package user

import "time"

type User struct {
    ID uint `json:"id" gorm:"primaryKey"`
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
	Email string `json:"email" gorm:"uniqueIndex;size:320"`
    Password string `json:"-"`
    Created time.Time `json:"created" gorm:"autoCreateTime"`
}

