package store

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	db "github.com/ariefzainuri96/go-logstream/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthStore struct {
	gormDb *db.GormDB
}

func (store *AuthStore) Login(ctx context.Context, body request.LoginRequest) (entity.User, string, error) {
	var user entity.User

	err := store.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		// get data by condition from user instance, which is by email
		return tx.
			Model(&entity.User{}).
			Where("email = ?", body.Email).
			// insert data to [user] address
			First(&user).Error
	})

	if err != nil {
		return user, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		return user, "", errors.New("invalid email or password")
	}

	token, err := generateToken(body.Email, int(user.ID))

	if err != nil {
		return user, "", err
	}

	return user, token, nil
}

func generateToken(email string, id int) (string, error) {
	jwtSecret := strings.TrimSpace(os.Getenv("SECRET_KEY"))

	claims := jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(), // Token valid for 30 day
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (store *AuthStore) Register(ctx context.Context, body request.RegisterRequest) (uint, error) {
	var emaiExists bool

	err := store.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.
			Model(&entity.User{}).
			Select("1"). // return 1 if email exists (this is signal that row exists)
			Where("email = ?", body.Email).
			Limit(1).          // stop query when row found
			Scan(&emaiExists). // the destination value is bool, and sql convert value from "1" to true
			Error
	})

	if err != nil {
		return 0, err
	}

	if emaiExists {
		return 0, errors.New("email sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		return 0, err
	}

	user := entity.User{
		Email:    body.Email,
		Password: string(hashedPassword),
	}

	err = store.gormDb.ExecWithTimeoutErr(ctx, func(tx *gorm.DB) error {
		return tx.Create(&user).Error
	})

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (store *AuthStore) ForgotPassword(ctx context.Context, body request.LoginRequest) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	result := store.gormDb.ExecWithTimeoutVal(ctx, func(tx *gorm.DB) *gorm.DB {
		return tx.
			Model(&entity.User{}).
			Where("email = ?", body.Email).
			Updates(entity.User{Password: string(hashedPassword)})
	})

	if result.Error != nil {
		return "", err
	} else if result.RowsAffected == 0 {
		return "", errors.New("email not found")
	}

	return "Success forgot password", nil
}
