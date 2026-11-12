package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/model"
	"inotal-be/internal/service"
)

type stubAboutCardStore struct {
	cards     []model.AboutCard
	nextID    uint
	createErr error
	listErr   error
}

func (s *stubAboutCardStore) List() ([]model.AboutCard, error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	return s.cards, nil
}

func (s *stubAboutCardStore) FindByID(id uint) (*model.AboutCard, error) {
	for i := range s.cards {
		if s.cards[i].ID == id {
			card := s.cards[i]
			return &card, nil
		}
	}
	return nil, nil
}

func (s *stubAboutCardStore) Create(card *model.AboutCard) error {
	if s.createErr != nil {
		return s.createErr
	}
	s.nextID++
	card.ID = s.nextID
	s.cards = append(s.cards, *card)
	return nil
}

func (s *stubAboutCardStore) Update(card *model.AboutCard) error {
	for i := range s.cards {
		if s.cards[i].ID == card.ID {
			s.cards[i] = *card
			return nil
		}
	}
	return service.ErrAboutCardNotFound
}

func (s *stubAboutCardStore) Delete(id uint) error {
	for i := range s.cards {
		if s.cards[i].ID == id {
			s.cards = append(s.cards[:i], s.cards[i+1:]...)
			return nil
		}
	}
	return service.ErrAboutCardNotFound
}

func newAboutCardTestRouter(store *stubAboutCardStore) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewAboutCardHandler(service.NewAboutCardService(store))

	router.GET("/api/about/cards", handler.GetCards)
	router.GET("/api/public/about/cards", handler.GetPublicCards)
	router.POST("/api/about/cards", handler.CreateCard)
	router.PUT("/api/about/cards/:id", handler.UpdateCard)
	router.DELETE("/api/about/cards/:id", handler.DeleteCard)
	return router
}

func performAboutCardRequest(t *testing.T, router *gin.Engine, method, path string, body any) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()

	var payload []byte
	if body != nil {
		var err error
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
	}

	request := httptest.NewRequest(method, path, bytes.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	var decoded map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
	}
	return recorder, decoded
}

func handlerAboutCardPayload() map[string]any {
	return map[string]any{
		"title":       "Solusi Terukur",
		"description": "Kolaborasi dimulai dari pemetaan kebutuhan mitra.",
		"icon":        "lightbulb",
		"url":         "#program",
		"sort_order":  1,
	}
}

func TestCreateCardReturns201AndLocation(t *testing.T) {
	router := newAboutCardTestRouter(&stubAboutCardStore{})

	recorder, body := performAboutCardRequest(t, router, http.MethodPost, "/api/about/cards", handlerAboutCardPayload())

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201 (body: %s)", recorder.Code, recorder.Body.String())
	}
	if location := recorder.Header().Get("Location"); location != "/api/about/cards/1" {
		t.Errorf("Location = %q, want /api/about/cards/1", location)
	}
	if body["success"] != true {
		t.Errorf("success = %v, want true", body["success"])
	}
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("data bukan object: %T", body["data"])
	}
	if data["sort_order"] != float64(1) {
		t.Errorf("sort_order = %v, want 1", data["sort_order"])
	}
}

func TestCreateCardRejectsInvalidPayload(t *testing.T) {
	router := newAboutCardTestRouter(&stubAboutCardStore{})
	payload := handlerAboutCardPayload()
	payload["icon"] = "ikon-ngawur"

	recorder, body := performAboutCardRequest(t, router, http.MethodPost, "/api/about/cards", payload)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (body: %s)", recorder.Code, recorder.Body.String())
	}
	if body["success"] != false {
		t.Errorf("success = %v, want false", body["success"])
	}
}

func TestCreateCardPositionConflictReturns409(t *testing.T) {
	router := newAboutCardTestRouter(&stubAboutCardStore{createErr: service.ErrAboutCardPositionOutOfRange})

	recorder, body := performAboutCardRequest(t, router, http.MethodPost, "/api/about/cards", handlerAboutCardPayload())

	if recorder.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409 (body: %s)", recorder.Code, recorder.Body.String())
	}
	if body["success"] != false {
		t.Errorf("success = %v, want false", body["success"])
	}
}

func TestUpdateUnknownCardReturns404(t *testing.T) {
	router := newAboutCardTestRouter(&stubAboutCardStore{})

	recorder, _ := performAboutCardRequest(t, router, http.MethodPut, "/api/about/cards/99", handlerAboutCardPayload())

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404 (body: %s)", recorder.Code, recorder.Body.String())
	}
}

func TestInvalidIDReturns400(t *testing.T) {
	router := newAboutCardTestRouter(&stubAboutCardStore{})

	recorder, _ := performAboutCardRequest(t, router, http.MethodDelete, "/api/about/cards/abc", nil)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400 (body: %s)", recorder.Code, recorder.Body.String())
	}
}

func TestDeleteCardReturnsSuccessWithNullData(t *testing.T) {
	store := &stubAboutCardStore{cards: []model.AboutCard{{
		ID:          3,
		Title:       "Kartu",
		Description: "Deskripsi kartu yang cukup panjang.",
		Icon:        "globe",
		URL:         "#kontak",
		SortOrder:   1,
	}}}
	router := newAboutCardTestRouter(store)

	recorder, body := performAboutCardRequest(t, router, http.MethodDelete, "/api/about/cards/3", nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body: %s)", recorder.Code, recorder.Body.String())
	}
	if body["success"] != true {
		t.Errorf("success = %v, want true", body["success"])
	}
	if body["data"] != nil {
		t.Errorf("data = %v, want null", body["data"])
	}
	if len(store.cards) != 0 {
		t.Errorf("sisa kartu = %d, want 0", len(store.cards))
	}
}

func TestPublicEndpointReturnsEmptyArray(t *testing.T) {
	router := newAboutCardTestRouter(&stubAboutCardStore{})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/public/about/cards", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	var body struct {
		Success bool              `json:"success"`
		Data    []model.AboutCard `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Data == nil {
		t.Fatal("data = null, want []")
	}
	if len(body.Data) != 0 {
		t.Errorf("len(data) = %d, want 0", len(body.Data))
	}
}
