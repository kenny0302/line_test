package usecase

import (
	"context"
	"fmt"
	"main/internal/application/dto"
	"main/internal/domain/user"
)

// ListUsersUseCase handles listing all users
type ListUsersUseCase struct {
	userRepo user.Repository
}

// NewListUsersUseCase creates a new ListUsersUseCase
func NewListUsersUseCase(userRepo user.Repository) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepo: userRepo,
	}
}

// Execute retrieves all users
func (uc *ListUsersUseCase) Execute(ctx context.Context) ([]dto.UserDTO, error) {
	users, err := uc.userRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	// Convert to DTOs
	userDTOs := make([]dto.UserDTO, 0, len(users))
	for _, u := range users {
		userDTOs = append(userDTOs, dto.UserDTO{
			UserID:        u.ID().Value(),
			DisplayName:   u.DisplayName().Value(),
			PictureURL:    u.PictureURL().Value(),
			StatusMessage: u.StatusMessage().Value(),
			Language:      u.Language().Value(),
			CreatedAt:     u.CreatedAt(),
			UpdatedAt:     u.UpdatedAt(),
		})
	}

	return userDTOs, nil
}
