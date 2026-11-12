package model

import "time"

// AboutEmpoweringSingletonID mengunci satu-satunya baris konfigurasi.
// ID tetap dipakai supaya duplikasi ditolak database, bukan cuma oleh service.
const AboutEmpoweringSingletonID uint = 1

// AboutEmpowering — konfigurasi label dan judul bagian "Memberdayakan Bisnis"
// di halaman publik Tentang. Singleton: satu baris untuk seluruh situs.
type AboutEmpowering struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"              json:"id"`
	Label       string    `gorm:"type:varchar(60);not null"             json:"label"`
	Title       string    `gorm:"type:varchar(150);not null"            json:"title"`
	TitleAccent string    `gorm:"type:varchar(100);not null"            json:"title_accent"`
	CreatedAt   time.Time `                                             json:"created_at"`
	UpdatedAt   time.Time `                                             json:"updated_at"`
}

// TableName memakai nama tunggal karena tabel ini menyimpan satu konfigurasi,
// bukan kumpulan baris.
func (AboutEmpowering) TableName() string {
	return "about_empowering"
}
