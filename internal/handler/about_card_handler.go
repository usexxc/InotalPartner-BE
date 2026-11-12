package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type AboutCardHandler struct {
	cardSvc *service.AboutCardService
}

func NewAboutCardHandler(cardSvc *service.AboutCardService) *AboutCardHandler {
	return &AboutCardHandler{cardSvc: cardSvc}
}

// GetCards — GET /api/about/cards (admin/super_admin)
func (h *AboutCardHandler) GetCards(c *gin.Context) {
	cards, err := h.cardSvc.ListCards()
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar kartu Tentang",
		"data":    cards,
	})
}

// GetPublicCards — GET /api/public/about/cards (tanpa login)
func (h *AboutCardHandler) GetPublicCards(c *gin.Context) {
	h.GetCards(c)
}

// CreateCard — POST /api/about/cards (admin/super_admin)
func (h *AboutCardHandler) CreateCard(c *gin.Context) {
	var req service.AboutCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data kartu Tentang tidak valid.",
		})
		return
	}

	card, err := h.cardSvc.CreateCard(req)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.Header("Location", "/api/about/cards/"+strconv.FormatUint(uint64(card.ID), 10))
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Kartu Tentang berhasil dibuat",
		"data":    card,
	})
}

// UpdateCard — PUT /api/about/cards/:id (admin/super_admin)
func (h *AboutCardHandler) UpdateCard(c *gin.Context) {
	id, ok := parseAboutCardID(c)
	if !ok {
		return
	}

	var req service.AboutCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data kartu Tentang tidak valid.",
		})
		return
	}

	card, err := h.cardSvc.UpdateCard(id, req)
	if err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kartu Tentang berhasil diperbarui",
		"data":    card,
	})
}

// DeleteCard — DELETE /api/about/cards/:id (admin/super_admin)
func (h *AboutCardHandler) DeleteCard(c *gin.Context) {
	id, ok := parseAboutCardID(c)
	if !ok {
		return
	}

	if err := h.cardSvc.DeleteCard(id); err != nil {
		h.writeError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kartu Tentang berhasil dihapus",
		"data":    nil,
	})
}

// writeError memetakan jenis galat ke status HTTP tanpa membocorkan detail SQL.
func (h *AboutCardHandler) writeError(c *gin.Context, err error) {
	var validation *service.AboutCardValidationError
	switch {
	case errors.As(err, &validation):
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": validation.Error(),
		})
	case errors.Is(err, service.ErrAboutCardNotFound):
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Kartu Tentang tidak ditemukan",
		})
	case errors.Is(err, service.ErrAboutCardPositionOutOfRange):
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Posisi kartu di luar rentang koleksi. Muat ulang daftar lalu pilih posisi yang tersedia.",
		})
	default:
		log.Printf("[about-cards] %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Terjadi kesalahan pada server. Silakan coba lagi.",
		})
	}
}

func parseAboutCardID(c *gin.Context) (uint, bool) {
	value, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || value == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID kartu Tentang tidak valid",
		})
		return 0, false
	}
	return uint(value), true
}
