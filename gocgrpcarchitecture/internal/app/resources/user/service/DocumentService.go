package service

import (
	"context"
	"gocgrpcarchitecture/internal/app/resources/user/models"
)

type DocumentService interface {
	Create(ctx context.Context, document *models.Document) error
	GetByID(ctx context.Context, Id string) (*models.Document, error)
	Update(ctx context.Context, Id string, Title string) error
	Delete(ctx context.Context, Id string) error
}

/*
type DocumentService interface {
	StreamUpload(ctx context.Context, filename string, payloadStream io.Reader) (*models.Document, error)
	GetDetails(ctx context.Context, Id string) (*models.Document, error)
	UpdateMetadata(ctx context.Context, id string, filename string) (*models.Document, error)
	PurgeDocument(ctx context.Context, id string) error
}
*/
//type DocumentService interface {
//	Save(document *models.Document) error
//	FindById(ctx context.Context, Id string) (*models.Document, error)
//	Update(ctx context.Context, document *models.Document) error
//	Delete(ctx context.Context, Id string) error
//}
