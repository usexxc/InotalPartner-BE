package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"inotal-be/internal/service"
)

type AboutEmpoweringHandler struct {
	empoweringSvc *service.AboutEmpoweringService
}

func NewAboutEmpoweringHandler(empoweringSvc *service.AboutEmpoweringService) *AboutEmpoweringHandler {
	return &AboutEmpoweringHandler{empoweringSvc: empoweringSvc}
}

// GetAboutEmpowering — GET /api/about/empowering (admin & super_admin)
// data:null berarti konfigurasi belum pernah dibuat.
func (h *AboutEmpoweringHandler) GetAboutEmpowering(c *gin.Context) {
	record, err := h.empoweringSvc.GetAboutEmpowering()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil pengaturan Memberdayakan Bisnis",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil pengaturan Memberdayakan Bisnis",
		"data":    record,
	})
}

// GetPublicAboutEmpowering — GET /api/public/about/empowering (publik)
// data:null dipakai frontend untuk kembali ke teks bawaan.
func (h *AboutEmpoweringHandler) GetPublicAboutEmpowering(c *gin.Context) {
	record, err := h.empoweringSvc.GetAboutEmpowering()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil pengaturan Memberdayakan Bisnis",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil pengaturan Memberdayakan Bisnis",
		"data":    record,
	})
}

// CreateAboutEmpowering — POST /api/about/empowering (admin & super_admin)
// Konfigurasi bersifat singleton: pembuatan kedua ditolak 409.
func (h *AboutEmpoweringHandler) CreateAboutEmpowering(c *gin.Context) {
	var req service.SaveAboutEmpoweringRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data request tidak valid",
		})
		return
	}

	record, err := h.empoweringSvc.CreateAboutEmpowering(req)
	if err != nil {
		h.writeAboutEmpoweringError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pengaturan Memberdayakan Bisnis berhasil dibuat",
		"data":    record,
	})
}

// UpdateAboutEmpowering — PUT /api/about/empowering (admin & super_admin)
// Mereplace tiga field sekaligus; 404 kalau konfigurasi belum pernah dibuat.
func (h *AboutEmpoweringHandler) UpdateAboutEmpowering(c *gin.Context) {
	var req service.SaveAboutEmpoweringRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data request tidak valid",
		})
		return
	}

	record, err := h.empoweringSvc.UpdateAboutEmpowering(req)
	if err != nil {
		h.writeAboutEmpoweringError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pengaturan Memberdayakan Bisnis berhasil diperbarui",
		"data":    record,
	})
}

// DeleteAboutEmpowering — DELETE /api/about/empowering (admin & super_admin)
// 204 tanpa body; GET sesudahnya mengembalikan data:null. Kartu tidak ikut terhapus.
func (h *AboutEmpoweringHandler) DeleteAboutEmpowering(c *gin.Context) {
	if err := h.empoweringSvc.DeleteAboutEmpowering(); err != nil {
		h.writeAboutEmpoweringError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AboutEmpoweringHandler) writeAboutEmpoweringError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	message := "Gagal memproses pengaturan Memberdayakan Bisnis"

	var invalid *service.InvalidAboutEmpoweringError
	switch {
	case errors.As(err, &invalid):
		status, message = http.StatusBadRequest, invalid.Reason
	case errors.Is(err, service.ErrAboutEmpoweringNotFound):
		status, message = http.StatusNotFound, err.Error()
	case errors.Is(err, service.ErrAboutEmpoweringExists):
		status, message = http.StatusConflict, err.Error()
	}

	c.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}
