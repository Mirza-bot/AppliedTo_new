package user

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
	Password *string `json:"password"`
}

type UserPublicDto struct {
    BaseUserDto
}
