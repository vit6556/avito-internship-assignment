package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vit6556/avito-internship-assignment/internal/database"
	"github.com/vit6556/avito-internship-assignment/internal/entity"
)

type MerchRepository struct {
	db *pgxpool.Pool
}

func NewMerchRepository(db *pgxpool.Pool) *MerchRepository {
	return &MerchRepository{
		db: db,
	}
}

func (r *MerchRepository) GetItemByID(ctx context.Context, itemID int) (*entity.MerchItem, error) {
	var item entity.MerchItem
	err := r.db.QueryRow(ctx, "SELECT id, name, price FROM merch_items WHERE id = $1", itemID).
		Scan(&item.ID, &item.Name, &item.Price)

	if err != nil {
		log.Printf("failed to get merch by ID %q: %v", itemID, err)
		return nil, database.ErrMerchNotFound
	}

	return &item, nil
}

func (r *MerchRepository) GetItemByName(ctx context.Context, name string) (*entity.MerchItem, error) {
	var item entity.MerchItem
	err := r.db.QueryRow(ctx, "SELECT id, name, price FROM merch_items WHERE name = $1", name).
		Scan(&item.ID, &item.Name, &item.Price)

	if err != nil {
		log.Printf("failed to get merch by name %q: %v", name, err)
		return nil, database.ErrMerchNotFound
	}

	return &item, nil
}

func (r *MerchRepository) GetUserPurchases(ctx context.Context, userID int) ([]*entity.InventoryItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT m.name, SUM(p.amount) as total_quantity
		FROM purchases p
		JOIN merch_items m ON p.item_id = m.id
		WHERE p.employee_id = $1
		GROUP BY m.name
	`, userID)

	if err != nil {
		log.Printf("failed to get purchases for user %d: %v", userID, err)
		return nil, database.ErrDatabaseQueryFailed
	}
	defer rows.Close()

	inventory := make([]*entity.InventoryItem, 0)
	for rows.Next() {
		var item entity.InventoryItem
		err := rows.Scan(&item.Type, &item.Quantity)
		if err != nil {
			log.Printf("failed to scan purchase row for user %d: %v", userID, err)
			return nil, database.ErrDatabaseScanFailed
		}
		inventory = append(inventory, &item)
	}

	return inventory, nil
}

func (r *MerchRepository) BuyItem(ctx context.Context, userID int, itemID int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin transaction for user %d: %v", userID, err)
		return database.ErrDatabaseTransaction
	}
	defer tx.Rollback(ctx)

	var item entity.MerchItem
	err = tx.QueryRow(ctx, "SELECT id, name, price FROM merch_items WHERE id = $1", itemID).
		Scan(&item.ID, &item.Name, &item.Price)
	if err != nil {
		log.Printf("failed to get merch by ID %q: %v", itemID, err)
		return database.ErrMerchNotFound
	}

	var employeeBalance int
	err = tx.QueryRow(ctx, "SELECT balance FROM employees WHERE id = $1", userID).Scan(&employeeBalance)
	if err != nil {
		log.Printf("failed to get balance for user %d: %v", userID, err)
		return database.ErrEmployeeNotFound
	}

	if employeeBalance < item.Price {
		return database.ErrInsufficientFunds
	}

	_, err = tx.Exec(ctx, "UPDATE employees SET balance = balance - $1 WHERE id = $2", item.Price, userID)
	if err != nil {
		log.Printf("failed to insert purchase record for user %d: %v", userID, err)
		return database.ErrDatabaseInsertFailed
	}

	_, err = tx.Exec(ctx, "INSERT INTO purchases (employee_id, item_id, amount) VALUES ($1, $2, $3)", userID, item.ID, 1)
	if err != nil {
		log.Printf("failed to commit transaction for user %d: %v", userID, err)
		return database.ErrDatabaseTransaction
	}

	return tx.Commit(ctx)
}
