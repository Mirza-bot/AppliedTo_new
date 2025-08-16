package mappers

import (
	"appliedTo/dtos/user_dtos"
	"appliedTo/models"
	"appliedTo/utils"
)

// --- INPUT MAPPERS ---

func CreateModel(dto userdtos.UserCreateDto) models.User {
	return models.User{
		FirstName: dto.FirstName,
		LastName: dto.LastName,
		Email: dto.Email,	
		Password: dto.Password,
	}
}

func PatchModel(m *models.User, dto userdtos.UserPatchDto) {
	utils.Patch(&m.FirstName, dto.FirstName)
	utils.Patch(&m.LastName, dto.LastName)
	utils.Patch(&m.Email, dto.Email)
}


// --- OUTPUT MAPPERS ---

func MapModelToPublicDto(u models.User) userdtos.UserPublicDto {
	return userdtos.UserPublicDto{
		BaseUserDto: userdtos.BaseUserDto{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
		},
	}
}


