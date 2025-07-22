package userdtos

type UserCreateDto struct {
    BaseUserDto
    Password string `json:"password"`
}

