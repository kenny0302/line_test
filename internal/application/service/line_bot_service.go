package service

import (
	"context"
	"fmt"
	"main/internal/application/dto"
	"main/internal/application/usecase"
)

// LINEBotService coordinates business operations for the LINE bot
type LINEBotService struct {
	registerUserUC *usecase.RegisterUserUseCase
	saveMessageUC  *usecase.SaveMessageUseCase
	listUsersUC    *usecase.ListUsersUseCase
	messagingClient MessagingClient
}

// MessagingClient defines the interface for sending messages
type MessagingClient interface {
	PushMessage(ctx context.Context, userID string, message string) error
	ReplyMessage(ctx context.Context, replyToken string, message string) error
}

// NewLINEBotService creates a new LINEBotService
func NewLINEBotService(
	registerUserUC *usecase.RegisterUserUseCase,
	saveMessageUC *usecase.SaveMessageUseCase,
	listUsersUC *usecase.ListUsersUseCase,
	messagingClient MessagingClient,
) *LINEBotService {
	return &LINEBotService{
		registerUserUC:  registerUserUC,
		saveMessageUC:   saveMessageUC,
		listUsersUC:     listUsersUC,
		messagingClient: messagingClient,
	}
}

// ProcessWebhookEvent processes a LINE webhook event
func (s *LINEBotService) ProcessWebhookEvent(ctx context.Context, userID, displayName, pictureURL, messageContent, replyToken string) error {
	// Register/update user
	registerReq := dto.RegisterUserRequest{
		UserID:      userID,
		DisplayName: displayName,
		PictureURL:  pictureURL,
	}

	_, err := s.registerUserUC.Execute(ctx, registerReq)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	// Save message
	if messageContent != "" {
		saveReq := dto.SaveMessageRequest{
			UserID:  userID,
			Content: messageContent,
		}

		_, err = s.saveMessageUC.Execute(ctx, saveReq)
		if err != nil {
			return fmt.Errorf("failed to save message: %w", err)
		}

		// Reply with the same message
		if replyToken != "" {
			err = s.messagingClient.ReplyMessage(ctx, replyToken, messageContent)
			if err != nil {
				return fmt.Errorf("failed to reply message: %w", err)
			}
		}
	}

	return nil
}

// ListAllUsers retrieves all registered users
func (s *LINEBotService) ListAllUsers(ctx context.Context) ([]dto.UserDTO, error) {
	return s.listUsersUC.Execute(ctx)
}

// SendPushMessage sends a push message to a user
func (s *LINEBotService) SendPushMessage(ctx context.Context, userID, message string) error {
	return s.messagingClient.PushMessage(ctx, userID, message)
}
