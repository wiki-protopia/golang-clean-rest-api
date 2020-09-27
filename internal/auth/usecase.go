package auth

import (
	"context"
	"github.com/AleksK1NG/api-mc/internal/models"
	"github.com/google/uuid"
)

// User repo interface
type UseCase interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.UserUpdate) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}
