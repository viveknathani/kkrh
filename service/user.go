package service

import (
	"errors"
	"fmt"
	"log"
	"net/mail"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
	"github.com/viveknathani/kkrh/cache"
	"github.com/viveknathani/kkrh/entity"
	"golang.org/x/crypto/bcrypt"
)

const (
	ageOfToken = time.Hour * 24 * 2
)

// Signup creates a new user in the database if all the checks pass.
func (service *Service) Signup(u *entity.User) error {

	if u == nil {
		return ErrNilUser
	}

	_, err := mail.ParseAddress(u.Email)
	if err != nil {
		log.Print(err)
		return ErrInvalidEmailFormat
	}

	user, err := service.repo.GetUser(u.Email)
	if err != nil {
		return err
	}
	if user != nil {
		return ErrEmailExists
	}

	if !isValidPassword(string(u.Password)) {
		return ErrInvalidPasswordFormat
	}

	hash, err := bcrypt.GenerateFromPassword(u.Password, bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return errors.New("bcrypt error, check logs")
	}

	u.Password = hash

	err = service.repo.CreateUser(u)
	if err != nil {
		return err
	}
	return nil
}

// Login creates a new JWT and returns it if there is no error.
func (service *Service) Login(u *entity.User) (string, error) {

	if u == nil {
		return "", ErrNilUser
	}

	user, err := service.repo.GetUser(u.Email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidEmailPassword
	}

	err = bcrypt.CompareHashAndPassword(user.Password, u.Password)
	if err != nil {
		return "", ErrInvalidEmailPassword
	}

	return service.createToken(user.Id)
}

// createToken will create a new JWT with id as payload and an expiry time
func (service *Service) createToken(id string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(ageOfToken).Unix(),
	})

	return token.SignedString(service.jwtSecret)
}

// VerifyAndDecodeToken will get the payload we need if the token is valid.
func (service *Service) VerifyAndDecodeToken(token string) (string, error) {

	if service.isBlacklistedToken(token) {
		return "", ErrInvalidToken
	}

	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return service.jwtSecret, nil
	})

	if claims, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		return claims["id"].(string), nil
	}

	log.Print(err)
	return "", ErrInvalidToken
}

// isValidPassword does a linear time check for password format.
func isValidPassword(password string) bool {

	const minLength = 8
	length := 0

	hasNumber := false
	hasUppercase := false
	hasLowercase := false
	hasSpecial := false

	for _, c := range password {

		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUppercase = true
		case unicode.IsLower(c):
			hasLowercase = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		default:
			return false
		}

		length++
	}

	return length >= minLength && hasNumber && hasLowercase && hasUppercase && hasSpecial
}

// Logout will put the JWT in the cache which acts as a blacklist.
func (service *Service) Logout(token string) error {
	return service.blacklistToken(token)
}

func (service *Service) blacklistToken(token string) error {
	_, err := cache.Set(service.conn, token, []byte("true"))
	return err
}

func (service *Service) isBlacklistedToken(token string) bool {
	res, err := cache.Get(service.conn, token)
	if err != nil {
		log.Print(err)
		return false
	}
	return string(res) == "true"
}
