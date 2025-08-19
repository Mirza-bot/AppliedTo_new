package user

import (
	"appliedTo/internal/platform/patch"
)

// --- INPUT MAPPERS ---

func CreateModel(dto UserCreateDto) User {
	return User{
		FirstName: dto.FirstName,
		LastName: dto.LastName,
		Email: dto.Email,	
		Password: dto.Password,
	}
}

func PatchModel(m *User, dto UserPatchDto) {
	patch.Patch(&m.FirstName, dto.FirstName)
	patch.Patch(&m.LastName, dto.LastName)
	patch.Patch(&m.Email, dto.Email)
}


// --- OUTPUT MAPPERS ---

func MapModelToPublicDto(u User) UserPublicDto {
	return UserPublicDto{
		BaseUserDto: BaseUserDto{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
		},
	}
}


