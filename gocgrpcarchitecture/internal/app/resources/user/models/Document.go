package models

/*
golang grpc production ready client and server user document file upload using gorm mysql with user
permission read, edit, update, delete, download with Service and Repository, add file type and size check,
maximum allowed size to be dynamic or configurable

like to enforce Role-Based Access Control (RBAC) instead of basic
Would you like the maximum allowed size to be dynamic / configurable via
Do you want to save the determined MIME type or exact file size into your MySQL
integrate secure cryptographic hashing (SHA-256) into the database to check for file duplicates?
like to enforce Role-Based Access Control (RBAC) instead of basic
I want to add file type and size check
*/
/*
type Document struct {
	Id            int64     `gorm:"column:Id;primary_key;autoIncrement" json:"Id"`
	FileName      string    `gorm:"column:FileName;size:255;not null;" json:"FileName"`
	Description   string    `gorm:"column:Description;size:255;not null;" json:"Description"`
	FileType      string    `gorm:"column:FileType;size:255;not null;" json:"FileType"`
	FileSize      int64     `gorm:"column:FileSize;not null;" json:"FileSize"`
	FileUrl       string    `gorm:"column:FileUrl;size:255;not null;" json:"FileUrl"`
	FilePath      string    `gorm:"column:FilePath;size:255;not null;" json:"FilePath"`
	OwnerId       int64     `gorm:"column:OwnerId;not null;" json:"OwnerId"`
	OwnerUniqueId string    `gorm:"column:OwnerUniqueId;size:255;not null;" json:"OwnerUniqueId"`
	ToPublic      bool      `gorm:"column:ToPublic;not null;" json:"ToPublic"`
	ToPrivate     bool      `gorm:"column:ToPrivate;not null;" json:"ToPrivate"`
	CreatedAt     time.Time `gorm:"column:CreatedAt" json:"CreatedAt"`
	UpdatedAt     time.Time `gorm:"column:UpdatedAt" json:"UpdatedAt"`
}

type Permission struct {
	Id         uint      `gorm:"column:Id;primary_key;" json:"Id"`
	DocumentId string    `gorm:"column:DocumentId;not null;" json:"DocumentId"`
	UserId     string    `gorm:"column:UserId;not null;" json:"UserId"`
	CanRead    bool      `gorm:"column:CanRead;not null;" json:"CanRead"`
	CanEdit    bool      `gorm:"column:CanEdit;not null;" json:"CanEdit"`
	CanDelete  bool      `gorm:"column:CanDelete;not null;" json:"CanDelete"`
	CreatedAt  time.Time `gorm:"column:CreatedAt" json:"CreatedAt"`
}
*/
