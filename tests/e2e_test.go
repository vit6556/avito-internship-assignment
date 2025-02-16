package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vit6556/avito-internship-assignment/internal/app"
	"github.com/vit6556/avito-internship-assignment/internal/config"
)

func setupTestAPI(t *testing.T) (func(), string, error) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join("..", "migrations", "000001_create_tables.up.sql"), filepath.Join("..", "migrations", "000002_add_merch_items.up.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatal("failed to get connection string")
	}

	dbPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal("failed to create db pool")
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		log.Fatalf("failed to ping db: %s", err.Error())
	}

	cfg := config.LoadServerConfig()
	e := app.InitServer(cfg, dbPool)
	testServer := httptest.NewServer(e)

	teardown := func() {
		log.Println("Stopping PostgreSQL container and shutting down server...")
		testServer.Close()
		_ = pgContainer.Terminate(ctx)
		dbPool.Close()
	}

	return teardown, testServer.URL, nil
}

func getAuthToken(baseURL, username, password string) (string, error) {
	authData := map[string]string{"username": username, "password": password}
	body, _ := json.Marshal(authData)

	resp, err := http.Post(baseURL+"/api/auth", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Token, nil
}

func TestSendCoinAPI(t *testing.T) {
	teardown, baseURL, err := setupTestAPI(t)
	if err != nil {
		t.Fatalf("failed to setup test API: %v", err)
	}
	defer teardown()

	senderToken, err := getAuthToken(baseURL, "sender", "password")
	assert.NoError(t, err)
	_, err = getAuthToken(baseURL, "receiver", "password")
	assert.NoError(t, err)

	reqBody := map[string]interface{}{
		"toUser": "receiver",
		"amount": 50,
	}
	body, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", baseURL+"/api/sendCoin", bytes.NewReader(body))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+senderToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestBuyItemAPI(t *testing.T) {
	teardown, baseURL, err := setupTestAPI(t)
	if err != nil {
		t.Fatalf("failed to setup test API: %v", err)
	}
	defer teardown()

	token, err := getAuthToken(baseURL, "buyer", "password")
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", baseURL+"/api/buy/book", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
