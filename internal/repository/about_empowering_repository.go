package repository

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"inotal-be/internal/model"
)

type AboutEmpoweringRepository struct {
	db *gorm.DB
}

func NewAboutEmpoweringRepository(db *gorm.DB) *AboutEmpoweringRepository {
	return &AboutEmpoweringRepository{db: db}
}

// Find mengambil konfigurasi yang ada. Hasil (nil, nil) berarti belum pernah dibuat.
func (r *AboutEmpoweringRepository) Find() (*model.AboutEmpowering, error) {
	var record model.AboutEmpowering
	err := r.db.Order("id ASC").First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// Create menyimpan konfigurasi singleton. Nilai balik false berarti baris sudah
// ada; bentrokan ditolak di level database sehingga dua request bersamaan tidak
// bisa membuat dua konfigurasi.
func (r *AboutEmpoweringRepository) Create(record *model.AboutEmpowering) (bool, error) {
	record.ID = model.AboutEmpoweringSingletonID
	result := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(record)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}

// Update mereplace tiga field konfigurasi sekaligus. Nilai balik false berarti
// konfigurasi belum ada.
func (r *AboutEmpoweringRepository) Update(record *model.AboutEmpowering) (bool, error) {
	result := r.db.Model(&model.AboutEmpowering{}).
		Where("id = ?", record.ID).
		Updates(map[string]any{
			"label":        record.Label,
			"title":        record.Title,
			"title_accent": record.TitleAccent,
		})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}

// Delete menghapus konfigurasi. Nilai balik false berarti tidak ada yang dihapus.
func (r *AboutEmpoweringRepository) Delete() (bool, error) {
	result := r.db.Where("id = ?", model.AboutEmpoweringSingletonID).
		Delete(&model.AboutEmpowering{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected == 1, nil
}
