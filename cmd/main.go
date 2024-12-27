package main

import (
    "log"
    "os"

    "github.com/gin-gonic/gin"

    "github.com/torvictorvic/seek-v2/internal/config"
    "github.com/torvictorvic/seek-v2/internal/handler"
    "github.com/torvictorvic/seek-v2/internal/repository"
    "github.com/torvictorvic/seek-v2/internal/security"
    "github.com/torvictorvic/seek-v2/internal/service"

    "github.com/swaggo/files"
    "github.com/swaggo/gin-swagger"
    _ "github.com/torvictorvic/seek-v2/docs" 
)

// @title Sistema de Gestión de Candidatos
// @version 1.0
// @description Esta API gestiona candidatos en un proceso de reclutamiento
// @termsOfService http://swagger.io/terms/

// @contact.name Soporte
// @contact.url http://www.ccccc.com/support
// @contact.email soporte@ccccc.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api

func main() {
    // Load environment variables
    _ = os.Setenv("JWT_SECRET", "") // En prod usar .env

    db := config.ConnectDB()
    defer db.Close()

    // Start repository and service
    candidateRepo := repository.NewCandidateRepository(db)
    candidateService := service.NewCandidateService(candidateRepo)
    candidateHandler := handler.NewCandidateHandler(candidateService)

    r := gin.Default()

    // Endpoint to generate token
    r.POST("/login", handler.GenerateToken)

    // JWT protected routes
    auth := r.Group("/api", security.AuthMiddleware())

    auth.POST("/candidates", candidateHandler.CreateCandidate)
    auth.GET("/candidates/:id", candidateHandler.GetCandidateByID)
    auth.GET("/candidates", candidateHandler.GetAllCandidates)
    auth.PUT("/candidates/:id", candidateHandler.UpdateCandidate)
    auth.DELETE("/candidates/:id", candidateHandler.DeleteCandidate)

    // Routes Swagger UI
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    log.Println("Server run http://localhost:8080")
    r.Run(":8080")
}
