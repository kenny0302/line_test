package usecase

import (
	"context"
	"fmt"
	"main/internal/application/dto"
	"main/internal/domain/shared"
	"main/internal/domain/user"
	"time"
)

// RegisterUserUseCase handles user registration
type RegisterUserUseCase struct {
	userRepo    user.Repository
	lockService shared.LockService
}

// NewRegisterUserUseCase creates a new RegisterUserUseCase
func NewRegisterUserUseCase(
	userRepo user.Repository,
	lockService shared.LockService,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:    userRepo,
		lockService: lockService,
	}
}

// Execute registers a new user or updates an existing one
func (uc *RegisterUserUseCase) Execute(ctx context.Context, req dto.RegisterUserRequest) (*dto.UserDTO, error) {
	// Create value objects
	userID, err := user.NewUserID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	displayName, err := user.NewDisplayName(req.DisplayName)
	if err != nil {
		return nil, fmt.Errorf("invalid display name: %w", err)
	}

	pictureURL := user.NewPictureURL(req.PictureURL)
	statusMessage := user.NewStatusMessage(req.StatusMessage)
	language := user.NewLanguage(req.Language)

	// Acquire distributed lock for this user
	lockResource := "user:" + userID.Value()
	lock, err := uc.lockService.AcquireLockWithRetry(ctx, lockResource, 10*time.Second, 3, 100*time.Millisecond)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}
	defer lock.Release(ctx)

	// Check if user already exists
	exists, err := uc.userRepo.Exists(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}

	var u *user.User

	if exists {
		// Update existing user
		u, err = uc.userRepo.FindByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to find user: %w", err)
		}
		u.UpdateProfile(displayName, pictureURL)
	} else {
		// Create new user
		u = user.NewUser(userID, displayName, pictureURL, statusMessage, language)
	}

	// Save user
	if err := uc.userRepo.Save(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Convert to DTO
	return &dto.UserDTO{
		UserID:        u.ID().Value(),
		DisplayName:   u.DisplayName().Value(),
		PictureURL:    u.PictureURL().Value(),
		StatusMessage: u.StatusMessage().Value(),
		Language:      u.Language().Value(),
		CreatedAt:     u.CreatedAt(),
		UpdatedAt:     u.UpdatedAt(),
	}, nil
}
