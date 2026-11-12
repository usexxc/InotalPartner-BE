package model

import "time"

// AboutCard — kartu "Pilih yang Terbaik" pada halaman publik Tentang.
// Posisi tampil disimpan di SortOrder dan selalu dijaga kontigu 1..N oleh repository.
type AboutCard struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"    json:"id"`
	Title       string    `gorm:"type:varchar(150);not null"  json:"title"`
	Description string    `gorm:"type:text;not null"          json:"description"`
	Icon        string    `gorm:"type:varchar(32);not null"   json:"icon"`
	URL         string    `gorm:"type:varchar(2048);not null" json:"url"`
	SortOrder   int       `gorm:"not null;default:1;index"    json:"sort_order"`
	CreatedAt   time.Time `                                    json:"created_at"`
	UpdatedAt   time.Time `                                    json:"updated_at"`
}
