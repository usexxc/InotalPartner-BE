package service

import (
	"fmt"
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"

	"inotal-be/internal/model"
	"inotal-be/internal/repository"
)

// Batas panjang mengikuti src/contracts/aboutCard.schema.json di frontend.
const (
	aboutCardTitleMin       = 2
	aboutCardTitleMax       = 150
	aboutCardDescriptionMin = 10
	aboutCardDescriptionMax = 2000
	aboutCardURLMax         = 2048
)

// Pesan galat dari repository diteruskan agar handler bisa memetakan status HTTP.
var (
	ErrAboutCardNotFound           = repository.ErrAboutCardNotFound
	ErrAboutCardPositionOutOfRange = repository.ErrAboutCardPosition
)

// AboutCardValidationError menandai input yang tidak lolos validasi (HTTP 400).
type AboutCardValidationError struct {
	Message string
}

func (e *AboutCardValidationError) Error() string { return e.Message }

func invalidAboutCard(message string) error {
	return &AboutCardValidationError{Message: message}
}

// Katalog ikon sama persis dengan src/data/aboutCardIcons.js.
var aboutCardIcons = map[string]struct{}{
	"lightbulb":    {},
	"award":        {},
	"headset":      {},
	"handshake":    {},
	"target":       {},
	"rocket":       {},
	"shield-check": {},
	"users":        {},
	"trending-up":  {},
	"code":         {},
	"briefcase":    {},
	"globe":        {},
}

var aboutCardHashRoutes = map[string]struct{}{
	"#landing":     {},
	"#tentang":     {},
	"#program":     {},
	"#event":       {},
	"#berita":      {},
	"#artikel":     {},
	"#partnership": {},
	"#kontak":      {},
}

// AboutCardRequest adalah payload tambah/edit; ID hanya dikirim lewat path.
type AboutCardRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	URL         string `json:"url"`
	SortOrder   int    `json:"sort_order"`
}

// AboutCardStore adalah bagian repository yang dipakai service.
type AboutCardStore interface {
	List() ([]model.AboutCard, error)
	FindByID(id uint) (*model.AboutCard, error)
	Create(card *model.AboutCard) error
	Update(card *model.AboutCard) error
	Delete(id uint) error
}

type AboutCardService struct {
	cards AboutCardStore
}

func NewAboutCardService(cards AboutCardStore) *AboutCardService {
	return &AboutCardService{cards: cards}
}

// ListCards mengirim seluruh kartu untuk dashboard dan halaman publik.
func (s *AboutCardService) ListCards() ([]model.AboutCard, error) {
	cards, err := s.cards.List()
	if err != nil {
		return nil, err
	}
	if cards == nil {
		cards = []model.AboutCard{}
	}
	return cards, nil
}

func (s *AboutCardService) CreateCard(req AboutCardRequest) (*model.AboutCard, error) {
	card, err := buildAboutCard(0, req)
	if err != nil {
		return nil, err
	}
	if err := s.cards.Create(card); err != nil {
		return nil, err
	}
	saved, err := s.cards.FindByID(card.ID)
	if err != nil {
		return nil, err
	}
	if saved == nil {
		return nil, ErrAboutCardNotFound
	}
	return saved, nil
}

func (s *AboutCardService) UpdateCard(id uint, req AboutCardRequest) (*model.AboutCard, error) {
	card, err := buildAboutCard(id, req)
	if err != nil {
		return nil, err
	}
	if err := s.cards.Update(card); err != nil {
		return nil, err
	}
	saved, err := s.cards.FindByID(id)
	if err != nil {
		return nil, err
	}
	if saved == nil {
		return nil, ErrAboutCardNotFound
	}
	return saved, nil
}

func (s *AboutCardService) DeleteCard(id uint) error {
	return s.cards.Delete(id)
}

// buildAboutCard memvalidasi payload lalu menyusun model. ID 0 berarti record baru.
func buildAboutCard(id uint, req AboutCardRequest) (*model.AboutCard, error) {
	title := strings.TrimSpace(req.Title)
	switch length := utf8.RuneCountInString(title); {
	case length < aboutCardTitleMin:
		return nil, invalidAboutCard("Judul kartu wajib diisi minimal 2 karakter.")
	case length > aboutCardTitleMax:
		return nil, invalidAboutCard(fmt.Sprintf("Judul kartu maksimal 150 karakter (saat ini %d).", length))
	}

	description := strings.TrimSpace(req.Description)
	switch length := utf8.RuneCountInString(description); {
	case length < aboutCardDescriptionMin:
		return nil, invalidAboutCard("Deskripsi kartu wajib diisi minimal 10 karakter.")
	case length > aboutCardDescriptionMax:
		return nil, invalidAboutCard(fmt.Sprintf("Deskripsi kartu maksimal 2000 karakter (saat ini %d).", length))
	}

	icon := strings.TrimSpace(req.Icon)
	if _, ok := aboutCardIcons[icon]; !ok {
		return nil, invalidAboutCard("Ikon kartu tidak dikenali. Pilih ikon dari katalog yang tersedia.")
	}

	target, err := validateAboutCardURL(req.URL)
	if err != nil {
		return nil, err
	}

	if req.SortOrder < 1 {
		return nil, invalidAboutCard("Posisi kartu minimal 1.")
	}

	return &model.AboutCard{
		ID:          id,
		Title:       title,
		Description: description,
		Icon:        icon,
		URL:         target,
		SortOrder:   req.SortOrder,
	}, nil
}

// validateAboutCardURL menerima URL http(s) absolut atau route publik berbentuk hash.
func validateAboutCardURL(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", invalidAboutCard("URL kartu wajib diisi.")
	}
	if utf8.RuneCountInString(value) > aboutCardURLMax {
		return "", invalidAboutCard("URL kartu maksimal 2048 karakter.")
	}

	for _, r := range value {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			return "", invalidAboutCard("URL kartu tidak boleh mengandung spasi atau karakter kontrol.")
		}
	}
	if strings.ContainsRune(value, '\\') {
		return "", invalidAboutCard("URL kartu tidak boleh mengandung backslash.")
	}

	if strings.HasPrefix(value, "/#") || strings.HasPrefix(value, "#") {
		if isAboutCardHashTarget(value) {
			return value, nil
		}
		return "", invalidAboutCard("Route kartu harus salah satu halaman publik, misalnya #program atau #berita/judul-berita.")
	}

	if strings.HasPrefix(value, "//") {
		return "", invalidAboutCard("URL kartu harus diawali http:// atau https://.")
	}

	parsed, err := url.Parse(value)
	if err != nil {
		return "", invalidAboutCard("URL kartu tidak valid.")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", invalidAboutCard("URL kartu harus diawali http:// atau https://, atau berupa route halaman publik.")
	}
	if parsed.Hostname() == "" {
		return "", invalidAboutCard("URL kartu harus memiliki alamat domain yang valid.")
	}
	if parsed.User != nil {
		return "", invalidAboutCard("URL kartu tidak boleh memuat username atau password.")
	}

	return value, nil
}

// isAboutCardHashTarget memeriksa route publik: #program, /#program, #berita/<slug>.
func isAboutCardHashTarget(value string) bool {
	route := strings.TrimPrefix(value, "/")
	if _, ok := aboutCardHashRoutes[route]; ok {
		return true
	}

	const newsPrefix = "#berita/"
	if !strings.HasPrefix(route, newsPrefix) {
		return false
	}
	slug := strings.TrimPrefix(route, newsPrefix)
	if slug == "" {
		return false
	}
	for _, r := range slug {
		isAllowed := r == '-' || r == '_' ||
			(r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
		if !isAllowed {
			return false
		}
	}
	return true
}
