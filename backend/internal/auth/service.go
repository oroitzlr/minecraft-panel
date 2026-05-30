package auth

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oroitz-lago-ramos/minecraft-panel/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	users *mongo.Collection
}

func NewService(db *mongo.Database) *Service {
	return &Service{
		users: db.Collection("users"),
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *Service) Login(username, password string) (string, error) {
	var user models.User
	err := s.users.FindOne(context.Background(), bson.M{
		"username": username,
	}).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errors.New("utilisateur introuvable")
		}
		return "", err
	}

	if !CheckPassword(password, user.Password) {
		return "", errors.New("mot de passe incorrect")
	}

	s.users.UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"lastLogin": time.Now()}},
	)

	return generateJWT(user.Username, user.Role)
}

func generateJWT(username, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET non défini")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	})

	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("algorithme invalide")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("token invalide")
	}

	return token.Claims.(jwt.MapClaims), nil
}

func (s *Service) CreateUser(username, password, role string) error {
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.users.InsertOne(context.Background(), models.User{
		Username:  username,
		Password:  hash,
		Role:      role,
		CreatedAt: time.Now(),
		LastLogin: time.Now(),
	})
	return err
}

func (s *Service) EnsureAdminExists() error {
	count, err := s.users.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	if count == 0 {
		password := os.Getenv("ADMIN_PASSWORD")
		if password == "" {
			password = "admin123"
		}
		log.Println("⚠️ Création de l'admin par défaut...")
		return s.CreateUser("admin", password, "admin")
	}
	return nil
}