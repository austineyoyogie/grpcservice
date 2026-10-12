package implementation

import (
	"context"
	"gocgrpcarchitecture/internal/app/resources/user/models"
	"gocgrpcarchitecture/internal/app/resources/user/service"

	"gorm.io/gorm"
)

type DocumentServiceImpl struct {
	D *gorm.DB
}

func NewDocumentServiceImpl(D *gorm.DB) service.DocumentService {
	return &DocumentServiceImpl{D: D}
}

func (D DocumentServiceImpl) Create(ctx context.Context, document *models.Document) error {
	tx := D.D.Begin()
	if err := D.D.WithContext(ctx).Model(&models.Document{}).Create(document).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (D DocumentServiceImpl) GetByID(ctx context.Context, Id string) (*models.Document, error) {
	tx := D.D.Begin()
	document := &models.Document{}
	if err := D.D.WithContext(ctx).Model(&models.Document{}).Where("Id = ?", Id).Find(document).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return document, tx.Commit().Error
}

func (D DocumentServiceImpl) Update(ctx context.Context, Id string, Title string) error {
	return D.D.WithContext(ctx).Model(&models.Document{}).Where("Id = ?", Id).Update("Title", Title).Error
}

func (D DocumentServiceImpl) Delete(ctx context.Context, Id string) error {
	return D.D.WithContext(ctx).Model(&models.Document{}).Where("Id = ?", Id).Delete(&models.Document{}).Error
}
