package repository

import (
	"context"
	"gocgrpcarchitecture/internal/app/resources/user/models"
)

type DocumentRepository interface {
	WithTransaction(ctx context.Context, fn func(txRepo DocumentRepository) error) error
	Create(ctx context.Context, doc *models.Document) error
	FindByID(ctx context.Context, id string) (*models.Document, error)
	FindByIDForUpdate(ctx context.Context, id string) (*models.Document, error)
	Update(ctx context.Context, doc *models.Document) error
	Delete(ctx context.Context, id string) error
}
