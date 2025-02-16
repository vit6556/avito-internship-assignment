package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vit6556/avito-internship-assignment/internal/database"
	"github.com/vit6556/avito-internship-assignment/internal/entity"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

func NewTransaction(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (r *TransactionRepository) GetCoinHistory(ctx context.Context, userID int) (*entity.CoinHistory, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			t.sender_id, sender.username AS sender_name,
			t.receiver_id, receiver.username AS receiver_name,
			t.amount
		FROM transactions t
		LEFT JOIN employees sender ON t.sender_id = sender.id
		LEFT JOIN employees receiver ON t.receiver_id = receiver.id
		WHERE t.sender_id = $1 OR t.receiver_id = $1
	`, userID)

	if err != nil {
		return &entity.CoinHistory{}, err
	}
	defer rows.Close()

	receivedMap := make(map[string]int)
	sentMap := make(map[string]int)
	for rows.Next() {
		var senderID, receiverID, amount int
		var senderName, receiverName string
		err := rows.Scan(&senderID, &senderName, &receiverID, &receiverName, &amount)
		if err != nil {
			return &entity.CoinHistory{}, err
		}

		if senderID == userID {
			sentMap[receiverName] += amount
		} else {
			receivedMap[senderName] += amount
		}
	}

	received := make([]entity.CoinTransaction, 0, len(receivedMap))
	for user, amount := range receivedMap {
		received = append(received, entity.CoinTransaction{
			User:   user,
			Amount: amount,
		})
	}

	sent := make([]entity.CoinTransaction, 0, len(sentMap))
	for user, amount := range sentMap {
		sent = append(sent, entity.CoinTransaction{
			User:   user,
			Amount: amount,
		})
	}

	return &entity.CoinHistory{
		Received: received,
		Sent:     sent,
	}, nil
}

func (r *TransactionRepository) SendCoins(ctx context.Context, senderID, receiverID, amount int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("failed to begin transaction for sending coins from user %d to user %d: %v", senderID, receiverID, err)
		return database.ErrDatabaseTransaction
	}
	defer tx.Rollback(ctx)

	var senderBalance int
	err = tx.QueryRow(ctx, "SELECT balance FROM employees WHERE id = $1", senderID).Scan(&senderBalance)
	if err != nil {
		log.Printf("failed to get balance for sender %d: %v", senderID, err)
		return database.ErrEmployeeNotFound
	}

	if senderBalance < amount {
		return database.ErrInsufficientFunds
	}

	_, err = tx.Exec(ctx, "UPDATE employees SET balance = balance - $1 WHERE id = $2", amount, senderID)
	if err != nil {
		log.Printf("failed to update sender %d balance: %v", senderID, err)
		return database.ErrDatabaseUpdateFailed
	}

	_, err = tx.Exec(ctx, "UPDATE employees SET balance = balance + $1 WHERE id = $2", amount, receiverID)
	if err != nil {
		log.Printf("failed to update receiver %d balance: %v", receiverID, err)
		return database.ErrDatabaseUpdateFailed
	}

	_, err = tx.Exec(ctx, "INSERT INTO transactions (sender_id, receiver_id, amount) VALUES ($1, $2, $3)", senderID, receiverID, amount)
	if err != nil {
		log.Printf("failed to insert transaction record for sender %d -> receiver %d: %v", senderID, receiverID, err)
		return database.ErrDatabaseInsertFailed
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("failed to commit transaction for sender %d -> receiver %d: %v", senderID, receiverID, err)
		return database.ErrDatabaseTransaction
	}

	return nil
}
