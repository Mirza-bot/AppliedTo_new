package userdtos

type BaseUserDto struct {
    ID uint `json:"id" gorm:"primaryKey"`
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
    Email string `json:"email" gorm:"unique"`
    IsAccount bool `json:"isAccount"`
}
