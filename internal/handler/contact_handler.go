package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type ContactHandler struct {
	contactSvc *service.ContactService
}

func NewContactHandler(contactSvc *service.ContactService) *ContactHandler {
	return &ContactHandler{contactSvc: contactSvc}
}

// SubmitContact — POST /api/contacts (PUBLIK, tanpa login)
// Dipakai form "Tinggalkan Pesan" di halaman Hubungi Kami.
func (h *ContactHandler) SubmitContact(c *gin.Context) {
	var req service.SubmitContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	contact, err := h.contactSvc.SubmitContact(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pesan berhasil dikirim",
		"data":    contact,
	})
}

// GetContacts — GET /api/contacts?page=&limit= (butuh login: operator ke atas)
// Pesan terbaru selalu di paling atas.
func (h *ContactHandler) GetContacts(c *gin.Context) {
	page := 1
	limit := 100
	if v := c.Query("page"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			page = n
		}
	}
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}

	contacts, err := h.contactSvc.ListContacts(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil daftar pesan",
		"data":    contacts,
	})
}

// MarkContactRead — PATCH /api/contacts/:id/read
// Body: { "is_read": true } → tandai sudah dibaca (false = belum dibaca).
// Kalau is_read tidak dikirim, default true.
func (h *ContactHandler) MarkContactRead(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req service.MarkContactReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	isRead := true
	if req.IsRead != nil {
		isRead = *req.IsRead
	}

	contact, err := h.contactSvc.MarkRead(id, isRead)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "pesan tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status pesan berhasil diperbarui",
		"data":    contact,
	})
}

// DeleteContact — DELETE /api/contacts/:id
func (h *ContactHandler) DeleteContact(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.contactSvc.DeleteContact(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "pesan tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pesan berhasil dihapus",
	})
}
