package userdtos

type UserLoginDto struct {
    Email string `json:"email" gorm:"unique"`
    Password string `json:"password"`
}

