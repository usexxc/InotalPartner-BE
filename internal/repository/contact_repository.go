package repository

import (
	"errors"

	"gorm.io/gorm"

	"inotal-be/internal/model"
)

type ContactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

// Create menyimpan pesan baru dari form publik
func (r *ContactRepository) Create(contact *model.Contact) error {
	return r.db.Create(contact).Error
}

// FindAll mengambil daftar pesan.
// PESAN TERBARU SELALU DI PALING ATAS: ORDER BY created_at DESC
// (id DESC sebagai tie-breaker kalau created_at-nya sama persis).
func (r *ContactRepository) FindAll(limit, offset int) ([]model.Contact, error) {
	var contacts []model.Contact
	query := r.db.Order("created_at DESC, id DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&contacts).Error
	return contacts, err
}

// FindByID mencari pesan berdasarkan ID
func (r *ContactRepository) FindByID(id uint) (*model.Contact, error) {
	var contact model.Contact
	err := r.db.First(&contact, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &contact, err
}

// UpdateStatus menyimpan perubahan status baca (unread/read)
func (r *ContactRepository) UpdateStatus(contact *model.Contact) error {
	return r.db.Save(contact).Error
}

// Delete menghapus pesan berdasarkan ID
func (r *ContactRepository) Delete(id uint) error {
	return r.db.Delete(&model.Contact{}, id).Error
}
