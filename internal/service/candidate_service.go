package service

import (
    "fmt"

    "github.com/torvictorvic/seek-v2/internal/domain"
    "github.com/torvictorvic/seek-v2/internal/repository"
)

type CandidateService interface {
    CreateCandidate(candidate domain.Candidate) (int, error)
    GetCandidateByID(id int) (*domain.Candidate, error)
    GetAllCandidates() ([]domain.Candidate, error)
    UpdateCandidate(candidate domain.Candidate) error
    DeleteCandidate(id int) error
}

type candidateServiceImpl struct {
    repo repository.CandidateRepository
}

func NewCandidateService(repo repository.CandidateRepository) CandidateService {
    return &candidateServiceImpl{repo: repo}
}

func (s *candidateServiceImpl) CreateCandidate(candidate domain.Candidate) (int, error) {
    // Ejemplo de validaciones
    if candidate.Name == "" || candidate.Email == "" {
        return 0, fmt.Errorf("The fields 'Name' and 'Email' are required")
    }
    // Luego llama al repositorio
    return s.repo.Create(candidate)
}

func (s *candidateServiceImpl) GetCandidateByID(id int) (*domain.Candidate, error) {
    return s.repo.GetByID(id)
}

func (s *candidateServiceImpl) GetAllCandidates() ([]domain.Candidate, error) {
    return s.repo.GetAll()
}

func (s *candidateServiceImpl) UpdateCandidate(candidate domain.Candidate) error {
    return s.repo.Update(candidate)
}

func (s *candidateServiceImpl) DeleteCandidate(id int) error {
    return s.repo.Delete(id)
}
