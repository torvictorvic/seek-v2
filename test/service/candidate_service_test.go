package service_test

import (
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "github.com/torvictorvic/seek-v2/internal/domain"
    "github.com/torvictorvic/seek-v2/internal/service"
)

// mockCandidateRepo implementa CandidateRepository usando testify/mock
type mockCandidateRepo struct {
    mock.Mock
}

func (m *mockCandidateRepo) Create(candidate domain.Candidate) (int, error) {
    args := m.Called(candidate)
    return args.Int(0), args.Error(1)
}
func (m *mockCandidateRepo) GetByID(id int) (*domain.Candidate, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*domain.Candidate), args.Error(1)
}
func (m *mockCandidateRepo) GetAll() ([]domain.Candidate, error) {
    args := m.Called()
    return args.Get(0).([]domain.Candidate), args.Error(1)
}
func (m *mockCandidateRepo) Update(candidate domain.Candidate) error {
    args := m.Called(candidate)
    return args.Error(0)
}
func (m *mockCandidateRepo) Delete(id int) error {
    args := m.Called(id)
    return args.Error(0)
}


func TestCreateCandidate_Success(t *testing.T) {
    mockRepo := new(mockCandidateRepo)
    svc := service.NewCandidateService(mockRepo)

    input := domain.Candidate{
        Name:           "Jane Doe",
        Email:          "jane@example.com",
        Gender:         "female",
        SalaryExpected: 35000,
    }

    // Configuramos mock: al llamar Create() con este input, devolvemos (1, nil)
    mockRepo.On("Create", input).Return(1, nil)

    id, err := svc.CreateCandidate(input)
    assert.NoError(t, err)
    assert.Equal(t, 1, id)

    mockRepo.AssertExpectations(t)
}

func TestCreateCandidate_MissingData(t *testing.T) {
    mockRepo := new(mockCandidateRepo)
    svc := service.NewCandidateService(mockRepo)

    // Dejamos Name y Email vacíos para forzar validación
    input := domain.Candidate{
        Name:  "",
        Email: "",
    }

    id, err := svc.CreateCandidate(input)
    assert.Error(t, err)
    assert.Equal(t, 0, id)
    assert.Contains(t, err.Error(), "required")

    // El repositorio no debe invocarse si falla la validación
    mockRepo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestGetCandidateByID_Success(t *testing.T) {
    mockRepo := new(mockCandidateRepo)
    svc := service.NewCandidateService(mockRepo)

    fakeCandidate := &domain.Candidate{
        ID:    10,
        Name:  "User Ten",
        Email: "user10@example.com",
    }

    mockRepo.On("GetByID", 10).Return(fakeCandidate, nil)

    result, err := svc.GetCandidateByID(10)
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "User Ten", result.Name)

    mockRepo.AssertExpectations(t)
}

func TestGetCandidateByID_NotFound(t *testing.T) {
    mockRepo := new(mockCandidateRepo)
    svc := service.NewCandidateService(mockRepo)

    // Retornamos nil, nil => no encontrado
    mockRepo.On("GetByID", 99).Return(nil, nil)

    result, err := svc.GetCandidateByID(99)
    assert.NoError(t, err)
    assert.Nil(t, result)

    mockRepo.AssertExpectations(t)
}


func TestUpdateCandidate_Success(t *testing.T) {
    mockRepo := new(mockCandidateRepo)
    svc := service.NewCandidateService(mockRepo)

    candidate := domain.Candidate{
        ID:             5,
        Name:           "New Name",
        Email:          "new@example.com",
        SalaryExpected: 45000,
    }

    mockRepo.On("Update", candidate).Return(nil)

    err := svc.UpdateCandidate(candidate)
    assert.NoError(t, err)

    mockRepo.AssertExpectations(t)
}

func TestUpdateCandidate_RepoError(t *testing.T) {
    mockRepo := new(mockCandidateRepo)
    svc := service.NewCandidateService(mockRepo)

    candidate := domain.Candidate{
        ID:    5,
        Name:  "Fail Name",
        Email: "fail@example.com",
    }

    mockRepo.On("Update", candidate).Return(errors.New("db error"))

    err := svc.UpdateCandidate(candidate)
    assert.Error(t, err)
    assert.Equal(t, "db error", err.Error())

    mockRepo.AssertExpectations(t)
}


func TestDeleteCandidate_Success(t *testing.T) {
    mockRepo := new(mockCandidateRepo)
    svc := service.NewCandidateService(mockRepo)

    mockRepo.On("Delete", 10).Return(nil)

    err := svc.DeleteCandidate(10)
    assert.NoError(t, err)

    mockRepo.AssertExpectations(t)
}

func TestDeleteCandidate_RepoError(t *testing.T) {
    mockRepo := new(mockCandidateRepo)
    svc := service.NewCandidateService(mockRepo)

    mockRepo.On("Delete", 10).Return(errors.New("delete error"))

    err := svc.DeleteCandidate(10)
    assert.Error(t, err)
    assert.Equal(t, "delete error", err.Error())

    mockRepo.AssertExpectations(t)
}
