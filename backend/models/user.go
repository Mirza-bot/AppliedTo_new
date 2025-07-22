package models

import "time"

type User struct {
    ID uint `json:"id" gorm:"primaryKey"`
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
    Email string `json:"email" gorm:"unique"`
    Password string `json:"-"`
    Created time.Time `json:"created" gorm:"autoCreateTime"`
}

