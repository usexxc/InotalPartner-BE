package model

import (
	"time"

	"gorm.io/gorm"
)

type ContactStatus string

const (
	ContactStatusUnread ContactStatus = "unread"
	ContactStatusRead   ContactStatus = "read"
)

// Contact — pesan masuk dari form "Tinggalkan Pesan" di halaman publik Hubungi Kami.
// Pengirim umumnya belum login, jadi tabel ini tidak terhubung ke tabel users.
type Contact struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"           json:"id"`
	CreatedAt time.Time      `                                          json:"created_at"`
	UpdatedAt time.Time      `                                          json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                              json:"-"`

	Name    string        `gorm:"type:varchar(150);not null"                 json:"name"`
	Email   string        `gorm:"type:varchar(255);not null"                 json:"email"`
	Phone   string        `gorm:"type:varchar(50)"                           json:"phone"`
	Subject string        `gorm:"type:varchar(200);not null"                 json:"subject"`
	Message string        `gorm:"type:text;not null"                         json:"message"`
	Status  ContactStatus `gorm:"type:varchar(20);not null;default:'unread'" json:"status"`
}
