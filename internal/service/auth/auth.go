package authservice

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/vit6556/avito-internship-assignment/internal/database"
	"github.com/vit6556/avito-internship-assignment/internal/entity"
	"github.com/vit6556/avito-internship-assignment/internal/service"
)

type AuthService struct {
	employeeRepo       database.EmployeeRepository
	jwtSecret          string
	tokenTTL           time.Duration
	defaultUserBalance int
}

func NewAuthService(employeeRepo database.EmployeeRepository, jwtSecret string, tokenTTL time.Duration, defaultUserBalance int) *AuthService {
	return &AuthService{
		employeeRepo:       employeeRepo,
		jwtSecret:          jwtSecret,
		tokenTTL:           tokenTTL,
		defaultUserBalance: defaultUserBalance,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *AuthService) AuthorizeUser(ctx context.Context, username, password string) (string, error) {
	employee, err := s.employeeRepo.GetEmployeeByUsername(ctx, username)
	if err != nil {
		passwordHash, err := hashPassword(password)
		if err != nil {
			log.Println("failed to hash password:", err)
			return "", service.ErrEmployeeCreationFailed
		}

		newEmployeeID, err := s.employeeRepo.CreateEmployee(
			ctx,
			entity.Employee{
				Username:     username,
				PasswordHash: passwordHash,
				Balance:      s.defaultUserBalance,
			},
		)
		if err != nil {
			log.Println("failed to create user:", err)
			return "", service.ErrEmployeeCreationFailed
		}

		employee = &entity.Employee{
			ID:           newEmployeeID,
			Username:     username,
			PasswordHash: passwordHash,
			Balance:      s.defaultUserBalance,
		}
	} else if !checkPasswordHash(password, employee.PasswordHash) {
		return "", service.ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": employee.ID,
		"exp":     time.Now().Add(s.tokenTTL).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		log.Println("failed to create jwt token:", err)
		return "", service.ErrAuthenticationFailed
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return 0, service.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, service.ErrInvalidToken
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, service.ErrInvalidToken
	}

	return int(userID), nil
}
