package userdtos

type BaseUserDto struct {
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
    Email string `json:"email" gorm:"unique"`
}

type UserCreateDto struct {
    BaseUserDto
    Password string `json:"password"`
}

type UserPatchDto struct {
    FirstName *string `json:"firstName"`
    LastName *string `json:"lastName"`
    Email *string `json:"email" gorm:"unique"`
}

type UserPublicDto struct {
    BaseUserDto
}
