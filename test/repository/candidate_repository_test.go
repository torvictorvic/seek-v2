package repository_test

import (
    "regexp"
    "testing"
    "time"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"

    "github.com/torvictorvic/seek-v2/internal/domain"
    "github.com/torvictorvic/seek-v2/internal/repository"
)

func TestCreateCandidate(t *testing.T) {
    // 1) Creamos la DB simulada
    db, mock, err := sqlmock.New()
    assert.NoError(t, err, "Error al crear el mock de la DB")
    defer db.Close()

    // 2) Creamos el repositorio usando la DB simulada
    repo := repository.NewCandidateRepository(db)

    // 3) Candidato de prueba
    candidate := domain.Candidate{
        Name:           "Test User",
        Email:          "test.user@example.com",
        Gender:         "female",
        SalaryExpected: 35000.0,
    }

    // 4) Esperamos que se ejecute un INSERT con estos campos
    insertQuery := regexp.QuoteMeta("INSERT INTO candidates (name, email, gender, salary_expected) VALUES (?, ?, ?, ?)")

    mock.ExpectExec(insertQuery).
        WithArgs(candidate.Name, candidate.Email, candidate.Gender, candidate.SalaryExpected).
        WillReturnResult(sqlmock.NewResult(1, 1)) // Devuelve ID=1, 1 fila afectada

    // 5) Llamamos al método
    id, err := repo.Create(candidate)

    // 6) Verificamos
    assert.NoError(t, err, "No debe ocurrir error al crear candidato")
    assert.Equal(t, 1, id, "El ID devuelto debe ser 1")

    // 7) Validamos que las expectativas se cumplan
    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}

func TestGetByIDCandidate_Found(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := repository.NewCandidateRepository(db)

    // 1) Simulamos una fila con campos no nulos para created_at y updated_at
    selectQuery := regexp.QuoteMeta("SELECT id, name, email, gender, salary_expected, created_at, updated_at FROM candidates WHERE id = ?")

    // 2) Creamos filas simuladas
    now := time.Now()
    rows := sqlmock.NewRows([]string{
        "id", "name", "email", "gender", "salary_expected", "created_at", "updated_at",
    }).AddRow(
        2,
        "John Doe",
        "john.doe@example.com",
        "male",
        40000.0,
        now,
        now,
    )

    mock.ExpectQuery(selectQuery).
        WithArgs(2).
        WillReturnRows(rows)

    // 3) Llamamos al método
    candidate, err := repo.GetByID(2)

    // 4) Verificamos
    assert.NoError(t, err, "No debe ocurrir error al obtener candidato")
    assert.NotNil(t, candidate, "Debe retornarse un candidato")
    assert.Equal(t, 2, candidate.ID)
    assert.Equal(t, "John Doe", candidate.Name)
    assert.Equal(t, "john.doe@example.com", candidate.Email)

    // 5) Validamos mock
    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}

func TestGetByIDCandidate_NotFound(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := repository.NewCandidateRepository(db)

    selectQuery := regexp.QuoteMeta("SELECT id, name, email, gender, salary_expected, created_at, updated_at FROM candidates WHERE id = ?")

    // Simulamos que no hay filas devueltas
    rows := sqlmock.NewRows([]string{
        "id", "name", "email", "gender", "salary_expected", "created_at", "updated_at",
    })

    mock.ExpectQuery(selectQuery).
        WithArgs(99).
        WillReturnRows(rows)

    candidate, err := repo.GetByID(99)
    assert.NoError(t, err, "No debe haber error, pero no hay filas")
    assert.Nil(t, candidate, "Si no existe, se espera nil")

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}

func TestUpdateCandidate(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := repository.NewCandidateRepository(db)

    updateQuery := regexp.QuoteMeta("UPDATE candidates SET name = ?, email = ?, gender = ?, salary_expected = ? WHERE id = ?")

    candidate := domain.Candidate{
        ID:             1,
        Name:           "Updated Name",
        Email:          "updated@example.com",
        Gender:         "female",
        SalaryExpected: 38000.0,
    }

    // Simulamos 1 fila afectada
    mock.ExpectExec(updateQuery).
        WithArgs(candidate.Name, candidate.Email, candidate.Gender, candidate.SalaryExpected, candidate.ID).
        WillReturnResult(sqlmock.NewResult(0, 1))

    err = repo.Update(candidate)
    assert.NoError(t, err)

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}

func TestDeleteCandidate(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    defer db.Close()

    repo := repository.NewCandidateRepository(db)

    deleteQuery := regexp.QuoteMeta("DELETE FROM candidates WHERE id = ?")

    mock.ExpectExec(deleteQuery).
        WithArgs(10).
        WillReturnResult(sqlmock.NewResult(0, 1))

    err = repo.Delete(10)
    assert.NoError(t, err)

    err = mock.ExpectationsWereMet()
    assert.NoError(t, err)
}
