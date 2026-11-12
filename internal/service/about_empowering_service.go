package service

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// SaveAboutEmpoweringRequest — payload POST/PUT /api/about/empowering.
// Aturan panjang mengikuti kontrak frontend: diukur dalam code point,
// spasi di ujung dirapikan, karakter kontrol ditolak.
type SaveAboutEmpoweringRequest struct {
	Label       string `json:"label"`
	Title       string `json:"title"`
	TitleAccent string `json:"title_accent"`
}

var (
	// ErrAboutEmpoweringNotFound — konfigurasi belum pernah dibuat.
	ErrAboutEmpoweringNotFound = errors.New("pengaturan Memberdayakan Bisnis belum ada")
	// ErrAboutEmpoweringExists — konfigurasi singleton sudah ada; pakai PUT.
	ErrAboutEmpoweringExists = errors.New("pengaturan Memberdayakan Bisnis sudah ada")
)

// InvalidAboutEmpoweringError menandai isian yang tidak lolos validasi kontrak.
// Pesannya dipakai langsung sebagai message respons 400.
type InvalidAboutEmpoweringError struct {
	Reason string
}

func (e *InvalidAboutEmpoweringError) Error() string { return e.Reason }

type AboutEmpoweringService struct {
	empoweringRepo *repository.AboutEmpoweringRepository
}

func NewAboutEmpoweringService(empoweringRepo *repository.AboutEmpoweringRepository) *AboutEmpoweringService {
	return &AboutEmpoweringService{empoweringRepo: empoweringRepo}
}

// GetAboutEmpowering mengambil konfigurasi aktif. nil berarti belum dibuat.
func (s *AboutEmpoweringService) GetAboutEmpowering() (*model.AboutEmpowering, error) {
	return s.empoweringRepo.Find()
}

// CreateAboutEmpowering membuat konfigurasi singleton.
func (s *AboutEmpoweringService) CreateAboutEmpowering(req SaveAboutEmpoweringRequest) (*model.AboutEmpowering, error) {
	record, err := toAboutEmpoweringRecord(req)
	if err != nil {
		return nil, err
	}

	created, err := s.empoweringRepo.Create(record)
	if err != nil {
		return nil, err
	}
	if !created {
		return nil, ErrAboutEmpoweringExists
	}
	return record, nil
}

// UpdateAboutEmpowering mereplace tiga field konfigurasi sekaligus lalu membaca
// balik baris tersimpan supaya updated_at yang dikembalikan akurat.
func (s *AboutEmpoweringService) UpdateAboutEmpowering(req SaveAboutEmpoweringRequest) (*model.AboutEmpowering, error) {
	patch, err := toAboutEmpoweringRecord(req)
	if err != nil {
		return nil, err
	}

	current, err := s.empoweringRepo.Find()
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, ErrAboutEmpoweringNotFound
	}

	patch.ID = current.ID
	updated, err := s.empoweringRepo.Update(patch)
	if err != nil {
		return nil, err
	}
	if !updated {
		return nil, ErrAboutEmpoweringNotFound
	}

	return s.empoweringRepo.Find()
}

// DeleteAboutEmpowering menghapus konfigurasi; halaman publik kembali ke teks bawaan.
func (s *AboutEmpoweringService) DeleteAboutEmpowering() error {
	deleted, err := s.empoweringRepo.Delete()
	if err != nil {
		return err
	}
	if !deleted {
		return ErrAboutEmpoweringNotFound
	}
	return nil
}

func toAboutEmpoweringRecord(req SaveAboutEmpoweringRequest) (*model.AboutEmpowering, error) {
	label, err := normalizeAboutEmpoweringText(req.Label, "Label kecil", 1, 60)
	if err != nil {
		return nil, err
	}
	title, err := normalizeAboutEmpoweringText(req.Title, "Judul utama", 2, 150)
	if err != nil {
		return nil, err
	}
	titleAccent, err := normalizeAboutEmpoweringText(req.TitleAccent, "Judul aksen (warna biru)", 0, 100)
	if err != nil {
		return nil, err
	}

	return &model.AboutEmpowering{
		Label:       label,
		Title:       title,
		TitleAccent: titleAccent,
	}, nil
}

// normalizeAboutEmpoweringText menolak karakter kontrol pada nilai mentah,
// memangkas spasi di ujung, lalu memeriksa panjang dalam code point.
func normalizeAboutEmpoweringText(raw, fieldLabel string, min, max int) (string, error) {
	for _, r := range raw {
		if r < 0x20 || r == 0x7f {
			return "", &InvalidAboutEmpoweringError{
				Reason: fmt.Sprintf("%s tidak boleh memuat karakter kontrol atau baris baru", fieldLabel),
			}
		}
	}

	value := strings.TrimSpace(raw)
	length := utf8.RuneCountInString(value)
	if length < min {
		if length == 0 {
			return "", &InvalidAboutEmpoweringError{Reason: fmt.Sprintf("%s wajib diisi", fieldLabel)}
		}
		return "", &InvalidAboutEmpoweringError{
			Reason: fmt.Sprintf("%s minimal %d karakter", fieldLabel, min),
		}
	}
	if length > max {
		return "", &InvalidAboutEmpoweringError{
			Reason: fmt.Sprintf("%s maksimal %d karakter", fieldLabel, max),
		}
	}

	return value, nil
}
