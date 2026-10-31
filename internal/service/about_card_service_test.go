package service

import (
	"errors"
	"strings"
	"testing"

	"inotal-be/internal/model"
)

type fakeAboutCardStore struct {
	cards       []model.AboutCard
	nextID      uint
	listErr     error
	createErr   error
	updateErr   error
	deleteErr   error
	lastCreated *model.AboutCard
	lastUpdated *model.AboutCard
}

func newFakeAboutCardStore(cards ...model.AboutCard) *fakeAboutCardStore {
	store := &fakeAboutCardStore{cards: append([]model.AboutCard(nil), cards...)}
	for _, card := range cards {
		if card.ID > store.nextID {
			store.nextID = card.ID
		}
	}
	return store
}

func (s *fakeAboutCardStore) List() ([]model.AboutCard, error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	return append([]model.AboutCard(nil), s.cards...), nil
}

func (s *fakeAboutCardStore) FindByID(id uint) (*model.AboutCard, error) {
	for i := range s.cards {
		if s.cards[i].ID == id {
			card := s.cards[i]
			return &card, nil
		}
	}
	return nil, nil
}

func (s *fakeAboutCardStore) Create(card *model.AboutCard) error {
	if s.createErr != nil {
		return s.createErr
	}
	s.nextID++
	card.ID = s.nextID
	s.cards = append(s.cards, *card)
	s.lastCreated = card
	return nil
}

func (s *fakeAboutCardStore) Update(card *model.AboutCard) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	for i := range s.cards {
		if s.cards[i].ID == card.ID {
			s.cards[i] = *card
			s.lastUpdated = card
			return nil
		}
	}
	return ErrAboutCardNotFound
}

func (s *fakeAboutCardStore) Delete(id uint) error {
	if s.deleteErr != nil {
		return s.deleteErr
	}
	for i := range s.cards {
		if s.cards[i].ID == id {
			s.cards = append(s.cards[:i], s.cards[i+1:]...)
			return nil
		}
	}
	return ErrAboutCardNotFound
}

func validAboutCardRequest() AboutCardRequest {
	return AboutCardRequest{
		Title:       "Solusi Terukur",
		Description: "Kolaborasi dimulai dari pemetaan kebutuhan mitra.",
		Icon:        "lightbulb",
		URL:         "#program",
		SortOrder:   1,
	}
}

func TestCreateCardTrimsFieldsAndReturnsSavedRecord(t *testing.T) {
	store := newFakeAboutCardStore()
	svc := NewAboutCardService(store)

	req := validAboutCardRequest()
	req.Title = "  Solusi Terukur  "
	req.Description = "  Kolaborasi dimulai dari pemetaan kebutuhan mitra.  "
	req.Icon = " lightbulb "

	card, err := svc.CreateCard(req)
	if err != nil {
		t.Fatalf("CreateCard() error = %v", err)
	}
	if card.ID == 0 {
		t.Fatal("CreateCard() tidak mengembalikan ID")
	}
	if card.Title != "Solusi Terukur" {
		t.Errorf("Title = %q, want %q", card.Title, "Solusi Terukur")
	}
	if card.Description != "Kolaborasi dimulai dari pemetaan kebutuhan mitra." {
		t.Errorf("Description tidak di-trim: %q", card.Description)
	}
	if card.Icon != "lightbulb" {
		t.Errorf("Icon = %q, want %q", card.Icon, "lightbulb")
	}
}

func TestAboutCardFieldValidation(t *testing.T) {
	longDescription := strings.Repeat("a", aboutCardDescriptionMax+1)
	tests := []struct {
		name    string
		mutate  func(*AboutCardRequest)
		wantMsg string
	}{
		{"judul kosong", func(r *AboutCardRequest) { r.Title = "   " }, "minimal 2"},
		{"judul 1 karakter", func(r *AboutCardRequest) { r.Title = "A" }, "minimal 2"},
		{"judul 151 karakter", func(r *AboutCardRequest) { r.Title = strings.Repeat("a", 151) }, "maksimal 150"},
		{"deskripsi pendek", func(r *AboutCardRequest) { r.Description = "pendek" }, "minimal 10"},
		{"deskripsi terlalu panjang", func(r *AboutCardRequest) { r.Description = longDescription }, "maksimal 2000"},
		{"ikon tidak dikenal", func(r *AboutCardRequest) { r.Icon = "api-key" }, "Ikon kartu tidak dikenali"},
		{"posisi nol", func(r *AboutCardRequest) { r.SortOrder = 0 }, "minimal 1"},
		{"posisi negatif", func(r *AboutCardRequest) { r.SortOrder = -3 }, "minimal 1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validAboutCardRequest()
			tt.mutate(&req)

			_, err := NewAboutCardService(newFakeAboutCardStore()).CreateCard(req)
			if err == nil {
				t.Fatal("CreateCard() error = nil, want validation error")
			}
			var validation *AboutCardValidationError
			if !errors.As(err, &validation) {
				t.Fatalf("error type = %T, want *AboutCardValidationError", err)
			}
			if !strings.Contains(validation.Message, tt.wantMsg) {
				t.Errorf("message = %q, want mengandung %q", validation.Message, tt.wantMsg)
			}
		})
	}
}

func TestAboutCardTitleUsesRuneCount(t *testing.T) {
	req := validAboutCardRequest()
	req.Title = "🚀🚀"

	if _, err := NewAboutCardService(newFakeAboutCardStore()).CreateCard(req); err != nil {
		t.Fatalf("judul dua emoji seharusnya valid, error = %v", err)
	}

	req.Title = strings.Repeat("a", aboutCardTitleMax)
	if _, err := NewAboutCardService(newFakeAboutCardStore()).CreateCard(req); err != nil {
		t.Fatalf("judul 150 karakter seharusnya valid, error = %v", err)
	}
}

func TestAboutCardURLValidation(t *testing.T) {
	accepted := []string{
		"https://contoh.com/program",
		"http://contoh.com",
		"#program",
		"/#program",
		"#berita/juara-2026",
		"/#berita/a_b-1",
		"#kontak",
	}
	for _, value := range accepted {
		t.Run("terima "+value, func(t *testing.T) {
			req := validAboutCardRequest()
			req.URL = value

			card, err := NewAboutCardService(newFakeAboutCardStore()).CreateCard(req)
			if err != nil {
				t.Fatalf("URL %q seharusnya diterima, error = %v", value, err)
			}
			if card.URL != value {
				t.Errorf("URL = %q, want %q", card.URL, value)
			}
		})
	}

	rejected := []string{
		"",
		"javascript:alert(1)",
		"data:text/html;base64,PHNjcmlwdD4=",
		"file:///etc/passwd",
		"mailto:halo@contoh.com",
		"tel:+628123",
		"//evil.example",
		"#dashboard",
		"#login",
		"program",
		"/program",
		"https://user:rahasia@contoh.com",
		"https://con toh.com",
		"https://contoh.com\\evil",
		"https://contoh.com/\nadmin",
		"https://" + strings.Repeat("a", aboutCardURLMax),
	}
	for _, value := range rejected {
		t.Run("tolak "+value, func(t *testing.T) {
			req := validAboutCardRequest()
			req.URL = value

			_, err := NewAboutCardService(newFakeAboutCardStore()).CreateCard(req)
			if err == nil {
				t.Fatalf("URL %q seharusnya ditolak", value)
			}
			var validation *AboutCardValidationError
			if !errors.As(err, &validation) {
				t.Fatalf("error type = %T, want *AboutCardValidationError", err)
			}
		})
	}
}

func TestCreateCardPropagatesPositionConflict(t *testing.T) {
	store := newFakeAboutCardStore()
	store.createErr = ErrAboutCardPositionOutOfRange

	_, err := NewAboutCardService(store).CreateCard(validAboutCardRequest())
	if !errors.Is(err, ErrAboutCardPositionOutOfRange) {
		t.Fatalf("error = %v, want ErrAboutCardPositionOutOfRange", err)
	}
}

func TestUpdateCardUsesPathID(t *testing.T) {
	store := newFakeAboutCardStore(model.AboutCard{
		ID:          7,
		Title:       "Judul Lama",
		Description: "Deskripsi lama yang cukup panjang.",
		Icon:        "target",
		URL:         "#tentang",
		SortOrder:   1,
	})

	req := validAboutCardRequest()
	req.SortOrder = 2

	card, err := NewAboutCardService(store).UpdateCard(7, req)
	if err != nil {
		t.Fatalf("UpdateCard() error = %v", err)
	}
	if card.ID != 7 {
		t.Errorf("ID = %d, want 7", card.ID)
	}
	if store.lastUpdated == nil || store.lastUpdated.ID != 7 {
		t.Fatal("UpdateCard() tidak meneruskan ID dari path")
	}
	if store.lastUpdated.SortOrder != 2 {
		t.Errorf("SortOrder = %d, want 2", store.lastUpdated.SortOrder)
	}
}

func TestUpdateAndDeleteUnknownCardReturnNotFound(t *testing.T) {
	svc := NewAboutCardService(newFakeAboutCardStore())

	if _, err := svc.UpdateCard(99, validAboutCardRequest()); !errors.Is(err, ErrAboutCardNotFound) {
		t.Errorf("UpdateCard() error = %v, want ErrAboutCardNotFound", err)
	}
	if err := svc.DeleteCard(99); !errors.Is(err, ErrAboutCardNotFound) {
		t.Errorf("DeleteCard() error = %v, want ErrAboutCardNotFound", err)
	}
}

func TestListCardsReturnsEmptySlice(t *testing.T) {
	cards, err := NewAboutCardService(newFakeAboutCardStore()).ListCards()
	if err != nil {
		t.Fatalf("ListCards() error = %v", err)
	}
	if cards == nil {
		t.Fatal("ListCards() = nil, want []")
	}
}
