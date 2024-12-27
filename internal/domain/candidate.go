package domain

import "time"

type Candidate struct {
    ID             int       `json:"id"`
    Name           string    `json:"name"`
    Email          string    `json:"email"`
    Gender         string    `json:"gender"`
    SalaryExpected float64   `json:"salary_expected"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}
