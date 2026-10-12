package models

import (
	"database/sql"
	"gocgrpcarchitecture/internal/pkg/utils"
	"strconv"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id                    uint64       `gorm:"column:Id;primary_key;auto_increment" json:"Id"`
	UniqueId              string       `gorm:"column:UniqueId;size:100;not null;unique" json:"UniqueId"`
	UUId                  string       `gorm:"column:UUId;size:100;not null;unique" json:"UUId"`
	FirstName             string       `gorm:"column:FirstName;size:45;not null;" json:"FirstName"`
	LastName              string       `gorm:"column:LastName;size:45;not null;" json:"LastName"`
	Email                 string       `gorm:"column:Email;size:45;not null;unique" json:"Email"`
	Password              string       `gorm:"column:Password;size:255;not null;" json:"Password"`
	Phone                 string       `gorm:"column:Phone;size:45;not null;" json:"Phone"`
	Biograph              string       `gorm:"column:Biograph;size:255;not null;" json:"Biograph"`
	EmailToken            string       `gorm:"column:EmailToken;size:255;not null;" json:"EmailToken"`
	TOTPFactorSid         string       `gorm:"column:TOTPFactorSid;size:255;not null;" json:"TOTPFactorSid"`
	QrCodeSecret          string       `gorm:"column:QrCodeSecret;size:45;not null;" json:"QrCodeSecret"`
	QrCodeImageUri        string       `gorm:"column:QrCodeImageUri;size:255;not null;" json:"QrCodeImageUri"`
	QrCodeImage           string       `gorm:"column:QrCodeImage;size:255;not null;" json:"QrCodeImage"`
	LastLoginDeviceName   string       `gorm:"column:LastLoginDeviceName;size:45;not null;" json:"LastLoginDeviceName"`
	LastLoginIpAddress    string       `gorm:"column:LastLoginIpAddress;size:45;not null;" json:"LastLoginIpAddress"`
	FailedAttempts        int          `gorm:"column:FailedAttempts;default:0" json:"FailedAttempts"`
	UsingMfa              bool         `gorm:"column:UsingMfa;default:false" json:"UsingMfa"`
	EmailVerified         sql.NullBool `gorm:"column:EmailVerified;default:false" json:"EmailVerified"`
	SMSVerified           bool         `gorm:"column:SMSVerified;default:false" json:"SMSVerified"`
	TOTPEnabled           bool         `gorm:"column:TOTPEnabled;default:false" json:"TOTPEnabled"`
	AccountEnabled        sql.NullBool `gorm:"column:AccountEnabled;default:false" json:"AccountEnabled"`
	AccountDisabled       sql.NullBool `gorm:"column:AccountDisabled;default:false" json:"AccountDisabled"`
	AccountDisabledAt     string       `gorm:"column:AccountDisabledAt;not null" json:"AccountDisabledAt"`
	AccountNonExpired     bool         `gorm:"column:AccountNonExpired;default:false" json:"AccountNonExpired"`
	AccountNonLocked      bool         `gorm:"column:AccountNonLocked;default:false" json:"AccountNonLocked"`
	LastLoginTimeStamp    sql.NullTime `gorm:"column:LastLoginTimestamp" json:"LastLoginTimestamp"`
	LastFailedAttemptTime sql.NullTime `gorm:"column:LastFailedAttemptTime" json:"LastFailedAttemptTime"`
	PasswordLastChangedAt time.Time    `gorm:"column:PasswordLastChangedAt" json:"PasswordLastChangedAt"`
	Roles                 []*Role      `gorm:"many2many:user_roles" json:"Roles"`
	CreatedAt             time.Time    `gorm:"column:CreatedAt" json:"CreatedAt"`
	UpdatedAt             time.Time    `gorm:"column:UpdatedAt" json:"UpdatedAt"`
}

type Role struct {
	Id          uint64    `gorm:"column:Id;primary_key;auto_increment" json:"Id"`
	Name        string    `gorm:"column:Name;size:45;not null;" json:"Name"`
	Permissions string    `gorm:"column:Permissions;size:255;not null;" json:"Permissions"`
	Users       *[]User   `gorm:"many2many:user_roles" json:"User"`
	CreatedAt   time.Time `gorm:"column:CreatedAt" json:"CreatedAt"`
	UpdatedAt   time.Time `gorm:"column:UpdatedAt" json:"UpdatedAt"`
}

type UserRole struct {
	Id        uint64    `gorm:"column:Id;primary_key;auto_increment" json:"Id"`
	UserId    uint64    `gorm:"column:user_id;" json:"UserId"`
	RoleId    uint64    `gorm:"column:role_id;" json:"RoleId"`
	CreatedAt time.Time `gorm:"column:CreatedAt" json:"CreatedAt"`
	UpdatedAt time.Time `gorm:"column:UpdatedAt" json:"UpdatedAt"`
}

type Document struct {
	Id        string    `gorm:"primaryKey;type:varchar(45);not null"`
	Filename  string    `gorm:"type:varchar(255);not null;index:idx_filename"`
	Title     string    `gorm:"type:varchar(255);not null"`
	MimeType  string    `gorm:"type:varchar(100);not null"`
	Size      uint64    `gorm:"type:bigint;not null"`
	Data      []byte    `gorm:"type:longblob;not null"`
	CreatedAt time.Time `gorm:"column:CreatedAt" json:"CreatedAt"`
	UpdatedAt time.Time `gorm:"column:UpdatedAt" json:"UpdatedAt"`
}

func (U *User) BeforeCreate(CD *gorm.DB) error {
	StrId, _ := uuid.NewV7()
	U.UUId = utils.StrUniqueID(4) +
		strconv.FormatInt(time.Now().Unix(), 10) + StrId.String() +
		strconv.FormatInt(time.Now().Unix(), 10) + utils.StrUniqueID(4)
	U.UniqueId = utils.StrUniqueID(4) + "-" + utils.StrUniqueID(2) + "-" + utils.StrUniqueID(4)

	return nil
}

func AutoMigration(DB *gorm.DB) {
	D := DB.Debug().Migrator().AutoMigrate(&User{}, &Role{}, &UserRole{})
	//D := DB.Debug().Migrator().DropTable(&User{}, &UserRole{})
	if D != nil {
		return
	}
}
