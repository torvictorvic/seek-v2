package repository

import (
    "database/sql"
    "fmt"

    "github.com/torvictorvic/seek-v2/internal/domain"
)

type CandidateRepository interface {
    Create(candidate domain.Candidate) (int, error)
    GetByID(id int) (*domain.Candidate, error)
    GetAll() ([]domain.Candidate, error)
    Update(candidate domain.Candidate) error
    Delete(id int) error
}

type candidateRepositoryImpl struct {
    db *sql.DB
}

func NewCandidateRepository(db *sql.DB) CandidateRepository {
    return &candidateRepositoryImpl{db: db}
}

func (r *candidateRepositoryImpl) Create(candidate domain.Candidate) (int, error) {
    query := `INSERT INTO candidates (name, email, gender, salary_expected) VALUES (?, ?, ?, ?)`
    result, err := r.db.Exec(query, candidate.Name, candidate.Email, candidate.Gender, candidate.SalaryExpected)
    if err != nil {
        return 0, fmt.Errorf("Error creating candidate: %w", err)
    }
    insertID, _ := result.LastInsertId()
    return int(insertID), nil
}

func (r *candidateRepositoryImpl) GetByID(id int) (*domain.Candidate, error) {
    query := `SELECT id, name, email, gender, salary_expected, created_at, updated_at FROM candidates WHERE id = ?`
    row := r.db.QueryRow(query, id)

    var c domain.Candidate
    err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Gender, &c.SalaryExpected, &c.CreatedAt, &c.UpdatedAt)
    if err == sql.ErrNoRows {
        return nil, nil
    } else if err != nil {
        return nil, fmt.Errorf("Error getting candidate by ID: %w", err)
    }
    return &c, nil
}

func (r *candidateRepositoryImpl) GetAll() ([]domain.Candidate, error) {
    query := `SELECT id, name, email, gender, salary_expected, created_at, updated_at FROM candidates`
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("Error getting candidate list: %w", err)
    }
    defer rows.Close()

    var candidates []domain.Candidate
    for rows.Next() {
        var c domain.Candidate
        if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Gender, &c.SalaryExpected, &c.CreatedAt, &c.UpdatedAt); err != nil {
            return nil, err
        }
        candidates = append(candidates, c)
    }
    return candidates, nil
}

func (r *candidateRepositoryImpl) Update(candidate domain.Candidate) error {
    query := `UPDATE candidates SET name = ?, email = ?, gender = ?, salary_expected = ? WHERE id = ?`
    _, err := r.db.Exec(query, candidate.Name, candidate.Email, candidate.Gender, candidate.SalaryExpected, candidate.ID)
    if err != nil {
        return fmt.Errorf("Error updating candidate: %w", err)
    }
    return nil
}

func (r *candidateRepositoryImpl) Delete(id int) error {
    query := `DELETE FROM candidates WHERE id = ?`
    _, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("Error deleting candidate: %w", err)
    }
    return nil
}
