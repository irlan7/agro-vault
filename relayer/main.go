package main

import (
"log"
"github.com/gin-gonic/gin"
"agrovault-backend/config"
"agrovault-backend/handlers"
"agrovault-backend/services"
)

func main() {
cfg := config.LoadConfig()
nearSvc := services.NewNearService(cfg)
h := handlers.NewHandler(nearSvc)

r := gin.Default()

// Serve Static Frontend
r.StaticFile("/", "./static/index.html")
r.Static("/static", "./static")

// Healthcheck
r.GET("/api/health", func(c *gin.Context) {
c.JSON(200, gin.H{
"status":   "healthy",
"contract": cfg.ContractID,
"network":  "testnet",
})
})

// Rute API
api := r.Group("/api/v1")
{
api.GET("/credit/:farmer_id", h.GetCreditStatus)
api.POST("/credit/spend", h.SpendCredit)
}

log.Printf("AgroVault Platform aktif di http://localhost:%s", cfg.Port)
if err := r.Run(":" + cfg.Port); err != nil {
log.Fatalf("Gagal menjalankan server: %v", err)
}
}
