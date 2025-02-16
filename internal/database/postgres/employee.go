package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vit6556/avito-internship-assignment/internal/database"
	"github.com/vit6556/avito-internship-assignment/internal/entity"
)

type EmployeeRepository struct {
	db *pgxpool.Pool
}

func NewEmployeeRepository(db *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{
		db: db,
	}
}

func (r *EmployeeRepository) GetEmployeeByUsername(ctx context.Context, username string) (*entity.Employee, error) {
	var employee entity.Employee
	err := r.db.QueryRow(ctx, "SELECT id, username, password_hash, balance FROM employees WHERE username = $1", username).
		Scan(&employee.ID, &employee.Username, &employee.PasswordHash, &employee.Balance)
	if err != nil {
		log.Printf("failed to get employee by username %q: %v", username, err)
		return nil, database.ErrEmployeeNotFound
	}

	return &employee, nil
}

func (r *EmployeeRepository) GetEmployeeByID(ctx context.Context, userID int) (*entity.Employee, error) {
	var employee entity.Employee
	err := r.db.QueryRow(ctx, "SELECT id, username, password_hash, balance FROM employees WHERE id = $1", userID).
		Scan(&employee.ID, &employee.Username, &employee.PasswordHash, &employee.Balance)
	if err != nil {
		log.Printf("failed to get employee by ID %q: %v", userID, err)
		return nil, database.ErrEmployeeNotFound
	}

	return &employee, nil
}

func (r *EmployeeRepository) CreateEmployee(ctx context.Context, employee entity.Employee) (int, error) {
	var userID int
	err := r.db.QueryRow(ctx, "INSERT INTO employees (username, password_hash, balance) VALUES ($1, $2, $3) RETURNING id",
		employee.Username, employee.PasswordHash, employee.Balance).Scan(&userID)

	if err != nil {
		log.Printf("failed to create employee %q: %v", employee.Username, err)
		return 0, database.ErrEmployeeCreationFailed
	}

	return userID, nil
}
