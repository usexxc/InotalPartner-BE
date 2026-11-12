package repository

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"inotal-be/internal/model"
)

// aboutCardLockKey mengunci seluruh koleksi kartu Tentang selama satu transaksi,
// sehingga dua permintaan tulis yang datang bersamaan tidak menghasilkan posisi ganda.
const aboutCardLockKey int64 = 1706636

var (
	// ErrAboutCardNotFound dipakai saat ID kartu tidak ada di database.
	ErrAboutCardNotFound = errors.New("kartu tentang tidak ditemukan")
	// ErrAboutCardPosition dipakai saat posisi di luar rentang koleksi.
	ErrAboutCardPosition = errors.New("posisi kartu tentang di luar rentang koleksi")
)

type AboutCardRepository struct {
	db *gorm.DB
}

func NewAboutCardRepository(db *gorm.DB) *AboutCardRepository {
	return &AboutCardRepository{db: db}
}

// List mengambil seluruh kartu, urut posisi lalu ID.
func (r *AboutCardRepository) List() ([]model.AboutCard, error) {
	var cards []model.AboutCard
	if err := r.db.Order("sort_order ASC, id ASC").Find(&cards).Error; err != nil {
		return nil, err
	}
	if cards == nil {
		cards = []model.AboutCard{}
	}
	return cards, nil
}

func (r *AboutCardRepository) FindByID(id uint) (*model.AboutCard, error) {
	var card model.AboutCard
	err := r.db.First(&card, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &card, nil
}

// Create menyisipkan kartu pada posisi card.SortOrder (1..N+1) lalu menggeser kartu di bawahnya.
func (r *AboutCardRepository) Create(card *model.AboutCard) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := lockAboutCards(tx); err != nil {
			return err
		}
		if err := compactAboutCards(tx); err != nil {
			return err
		}

		var total int64
		if err := tx.Model(&model.AboutCard{}).Count(&total).Error; err != nil {
			return err
		}
		if card.SortOrder < 1 || int64(card.SortOrder) > total+1 {
			return ErrAboutCardPosition
		}

		if err := shiftAboutCardRange(tx, card.SortOrder, total, 1); err != nil {
			return err
		}
		return tx.Create(card).Error
	})
}

// Update memindahkan kartu ke posisi baru (1..N) lalu menggeser kartu pada interval lama/baru.
func (r *AboutCardRepository) Update(card *model.AboutCard) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := lockAboutCards(tx); err != nil {
			return err
		}
		if err := compactAboutCards(tx); err != nil {
			return err
		}

		var current model.AboutCard
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&current, card.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAboutCardNotFound
		}
		if err != nil {
			return err
		}

		var total int64
		if err := tx.Model(&model.AboutCard{}).Count(&total).Error; err != nil {
			return err
		}
		if card.SortOrder < 1 || int64(card.SortOrder) > total {
			return ErrAboutCardPosition
		}

		switch {
		case card.SortOrder < current.SortOrder:
			if err := shiftAboutCardRange(tx, card.SortOrder, int64(current.SortOrder-1), 1); err != nil {
				return err
			}
		case card.SortOrder > current.SortOrder:
			if err := shiftAboutCardRange(tx, current.SortOrder+1, int64(card.SortOrder), -1); err != nil {
				return err
			}
		}

		return tx.Model(&model.AboutCard{}).Where("id = ?", card.ID).Updates(map[string]interface{}{
			"title":       card.Title,
			"description": card.Description,
			"icon":        card.Icon,
			"url":         card.URL,
			"sort_order":  card.SortOrder,
		}).Error
	})
}

// Delete menghapus kartu lalu merapatkan posisi kartu sesudahnya.
func (r *AboutCardRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := lockAboutCards(tx); err != nil {
			return err
		}
		if err := compactAboutCards(tx); err != nil {
			return err
		}

		var card model.AboutCard
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&card, id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAboutCardNotFound
		}
		if err != nil {
			return err
		}

		if err := tx.Delete(&model.AboutCard{}, id).Error; err != nil {
			return err
		}
		return tx.Model(&model.AboutCard{}).
			Where("sort_order > ?", card.SortOrder).
			UpdateColumn("sort_order", gorm.Expr("sort_order - ?", 1)).Error
	})
}

// lockAboutCards memakai advisory lock transaksi agar perubahan posisi diserialisasi.
func lockAboutCards(tx *gorm.DB) error {
	return tx.Exec("SELECT pg_advisory_xact_lock(?)", aboutCardLockKey).Error
}

// compactAboutCards merapikan posisi menjadi 1..N sesuai urutan (sort_order, id).
// Dipanggil di dalam transaksi yang sudah memegang advisory lock.
func compactAboutCards(tx *gorm.DB) error {
	var cards []model.AboutCard
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Order("sort_order ASC, id ASC").
		Find(&cards).Error; err != nil {
		return err
	}

	for i := range cards {
		want := i + 1
		if cards[i].SortOrder == want {
			continue
		}
		if err := tx.Model(&model.AboutCard{}).
			Where("id = ?", cards[i].ID).
			UpdateColumn("sort_order", want).Error; err != nil {
			return err
		}
	}
	return nil
}

// shiftAboutCardRange menggeser sort_order pada rentang from..to sebanyak delta.
func shiftAboutCardRange(tx *gorm.DB, from int, to int64, delta int) error {
	if int64(from) > to {
		return nil
	}
	return tx.Model(&model.AboutCard{}).
		Where("sort_order >= ? AND sort_order <= ?", from, to).
		UpdateColumn("sort_order", gorm.Expr("sort_order + ?", delta)).Error
}
