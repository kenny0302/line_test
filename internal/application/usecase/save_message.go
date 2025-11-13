package usecase

import (
	"context"
	"fmt"
	"main/internal/application/dto"
	"main/internal/domain/message"
	"main/internal/domain/user"
	"strconv"
	"time"
)

// SaveMessageUseCase handles saving messages
type SaveMessageUseCase struct {
	messageRepo message.Repository
}

// NewSaveMessageUseCase creates a new SaveMessageUseCase
func NewSaveMessageUseCase(messageRepo message.Repository) *SaveMessageUseCase {
	return &SaveMessageUseCase{
		messageRepo: messageRepo,
	}
}

// Execute saves a new message
func (uc *SaveMessageUseCase) Execute(ctx context.Context, req dto.SaveMessageRequest) (*dto.MessageDTO, error) {
	// Create value objects
	userID, err := user.NewUserID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	content, err := message.NewContent(req.Content)
	if err != nil {
		return nil, fmt.Errorf("invalid message content: %w", err)
	}

	// Generate message ID from timestamp
	timestamp := time.Now()
	messageIDStr := strconv.FormatInt(timestamp.UnixNano(), 10)
	messageID, err := message.NewMessageID(messageIDStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create message ID: %w", err)
	}

	userRef := message.NewUserReference(userID)

	// Create message entity
	msg := message.NewMessage(messageID, userRef, content)

	// Save message
	if err := uc.messageRepo.Save(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	// Convert to DTO
	return &dto.MessageDTO{
		MessageID: msg.ID().Value(),
		UserID:    msg.UserRef().UserID().Value(),
		Content:   msg.Content().Value(),
		Timestamp: msg.Timestamp(),
	}, nil
}
