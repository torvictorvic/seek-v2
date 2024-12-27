package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/torvictorvic/seek-v2/internal/domain"
    "github.com/torvictorvic/seek-v2/internal/service"
)

type CandidateHandler struct {
    service service.CandidateService
}

func NewCandidateHandler(s service.CandidateService) *CandidateHandler {
    return &CandidateHandler{service: s}
}

// CreateCandidate godoc
// @Summary Crear un nuevo candidato
// @Description Crea un candidato con los datos enviados en el body
// @Tags Candidates
// @Accept  json
// @Produce  json
// @Param candidate body domain.Candidate true "Datos del candidato"
// @Success 200 {object} map[string]interface{} "ok"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /candidates [post]
// @Security Bearer
func (h *CandidateHandler) CreateCandidate(c *gin.Context) {
    var candidate domain.Candidate
    if err := c.ShouldBindJSON(&candidate); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "JSON no valid"})
        return
    }

    id, err := h.service.CreateCandidate(candidate)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Candidate created", "id": id})
}

// GetCandidateByID godoc
// @Summary Obtener candidato por ID
// @Description Retorna el candidato cuyo ID se pasa como parámetro
// @Tags Candidates
// @Accept  json
// @Produce  json
// @Param  id path int true "ID del Candidato"
// @Success 200 {object} domain.Candidate
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} map[string]interface{} "Candidato no encontrado"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /candidates/{id} [get]
// @Security Bearer
func (h *CandidateHandler) GetCandidateByID(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "The ID must be an integer"})
        return
    }

    candidate, err := h.service.GetCandidateByID(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if candidate == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Candidate not found"})
        return
    }

    c.JSON(http.StatusOK, candidate)
}

// GetAllCandidates godoc
// @Summary Listar todos los candidatos
// @Description Retorna la lista de todos los candidatos
// @Tags Candidates
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.Candidate
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /candidates [get]
// @Security Bearer
func (h *CandidateHandler) GetAllCandidates(c *gin.Context) {
    candidates, err := h.service.GetAllCandidates()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, candidates)
}

// UpdateCandidate godoc
// @Summary Actualiza un candidato
// @Description Actualiza un candidato con los datos enviados en el body
// @Tags Candidates
// @Accept  json
// @Produce  json
// @Param  id path int true "ID del Candidato"
// @Success 200 {object} map[string]interface{} "ok"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /candidates/{id} [put]
// @Security Bearer
func (h *CandidateHandler) UpdateCandidate(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "The ID must be an integer"})
        return
    }

    var candidate domain.Candidate
    if err := c.ShouldBindJSON(&candidate); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "JSON no valid"})
        return
    }
    candidate.ID = id

    err = h.service.UpdateCandidate(candidate)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Updated Candidate"})
}

// DeleteCandidate godoc
// @Summary Borra un candidato por ID
// @Description Borra un candidato cuyo ID se pasa como parámetro
// @Tags Candidates
// @Accept  json
// @Produce  json
// @Param  id path int true "ID del Candidato"
// @Success 200 {object} domain.Candidate
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} map[string]interface{} "Candidato no encontrado"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /candidates/{id} [delete]
// @Security Bearer
func (h *CandidateHandler) DeleteCandidate(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "The ID must be an integer"})
        return
    }

    err = h.service.DeleteCandidate(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Candidate eliminated"})
}
