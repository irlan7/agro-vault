package handlers

import (
"net/http"
"github.com/gin-gonic/gin"
"agrovault-backend/services"
)

type Handler struct {
nearSvc *services.NearService
}

func NewHandler(nearSvc *services.NearService) *Handler {
return &Handler{nearSvc: nearSvc}
}

// Cek Saldo & Kuota Kredit Petani
func (h *Handler) GetCreditStatus(c *gin.Context) {
farmerID := c.Param("farmer_id")
if farmerID == "" {
c.JSON(http.StatusBadRequest, gin.H{"error": "farmer_id harus diisi"})
return
}

status, err := h.nearSvc.GetFarmerCredit(farmerID)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, gin.H{
"farmer_id": farmerID,
"data":      status,
})
}

type SpendRequest struct {
FarmerID    string `json:"farmer_id" binding:"required"`
VendorID    string `json:"vendor_id" binding:"required"`
AmountYocto string `json:"amount_yocto" binding:"required"`
ItemName    string `json:"item_name"`
}

// Eksekusi Pembelian Pupuk (Scan QR Kios)
func (h *Handler) SpendCredit(c *gin.Context) {
var req SpendRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
return
}

txResult, err := h.nearSvc.SpendCreditRelay(req.FarmerID, req.VendorID, req.AmountYocto)
if err != nil {
c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
return
}

c.JSON(http.StatusOK, gin.H{
"status":     "SUCCESS",
"message":    "Pembayaran pupuk ke kios berhasil diverifikasi di blockchain NEAR",
"item":       req.ItemName,
"raw_output": txResult,
})
}
