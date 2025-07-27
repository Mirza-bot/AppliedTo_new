package userdtos

type BaseUserDto struct {
    ID uint `json:"id" gorm:"primaryKey"`
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
    Email string `json:"email" gorm:"unique"`
}

type UserCreateDto struct {
    BaseUserDto
    Password string `json:"password"`
}

type UserLoginDto struct {
    Email string `json:"email" gorm:"unique"`
    Password string `json:"password"`
}

type UserModifyDto struct {
    BaseUserDto
}

type UserPublicDto struct {
    BaseUserDto
}
