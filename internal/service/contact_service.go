package service

import (
	"errors"
	"fmt"
	"strings"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// ─── DTOs ─────────────────────────────────────────────────────────────────────

// SubmitContactRequest — payload dari form "Tinggalkan Pesan" (publik).
type SubmitContactRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=150"`
	Email   string `json:"email" binding:"required,email,max=255"`
	Phone   string `json:"phone" binding:"omitempty,max=50"`
	Subject string `json:"subject" binding:"required,min=3,max=200"`
	Message string `json:"message" binding:"required,min=10"`
}

// MarkContactReadRequest — body PATCH /api/contacts/:id/read
// IsRead pointer supaya bisa dibedakan "tidak dikirim" vs "false".
type MarkContactReadRequest struct {
	IsRead *bool `json:"is_read"`
}

// ─── Service ──────────────────────────────────────────────────────────────────

type ContactService struct {
	contactRepo *repository.ContactRepository
}

func NewContactService(contactRepo *repository.ContactRepository) *ContactService {
	return &ContactService{contactRepo: contactRepo}
}

// SubmitContact menyimpan pesan baru dari form publik (status awal: unread)
func (s *ContactService) SubmitContact(req SubmitContactRequest) (*model.Contact, error) {
	contact := &model.Contact{
		Name:    strings.TrimSpace(req.Name),
		Email:   strings.TrimSpace(strings.ToLower(req.Email)),
		Phone:   strings.TrimSpace(req.Phone),
		Subject: strings.TrimSpace(req.Subject),
		Message: strings.TrimSpace(req.Message),
		Status:  model.ContactStatusUnread,
	}

	if contact.Name == "" || contact.Subject == "" || contact.Message == "" {
		return nil, errors.New("data pesan tidak lengkap")
	}

	if err := s.contactRepo.Create(contact); err != nil {
		return nil, fmt.Errorf("gagal menyimpan pesan: %w", err)
	}
	return contact, nil
}

// ListContacts mengambil daftar pesan — pesan terbaru di paling atas.
// page mulai dari 1, limit maksimal 100.
func (s *ContactService) ListContacts(page, limit int) ([]model.Contact, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	contacts, err := s.contactRepo.FindAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return contacts, nil
}

// MarkRead menandai pesan sudah dibaca (true) atau belum dibaca (false)
func (s *ContactService) MarkRead(id uint, isRead bool) (*model.Contact, error) {
	contact, err := s.contactRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	if contact == nil {
		return nil, errors.New("pesan tidak ditemukan")
	}

	if isRead {
		contact.Status = model.ContactStatusRead
	} else {
		contact.Status = model.ContactStatusUnread
	}

	if err := s.contactRepo.UpdateStatus(contact); err != nil {
		return nil, fmt.Errorf("gagal memperbarui status pesan: %w", err)
	}
	return contact, nil
}

// DeleteContact menghapus pesan berdasarkan ID
func (s *ContactService) DeleteContact(id uint) error {
	contact, err := s.contactRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if contact == nil {
		return errors.New("pesan tidak ditemukan")
	}
	return s.contactRepo.Delete(id)
}
