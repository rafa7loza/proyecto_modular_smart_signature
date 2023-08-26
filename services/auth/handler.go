package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DBConn                 *gorm.DB
	Domain, SendGridAPIKey string
	Options                *AuthOptions
}

type AuthOptions struct {
	IsEmailVerficationEnabled *bool
}

func NewAuthHandler(dsn string, domain string, apiKey string, opts *AuthOptions) *AuthHandler {
	db := connectDB(dsn)
	if opts.IsEmailVerficationEnabled == nil {
		// By default should be disabled
		isVerificationEnabled := false
		opts.IsEmailVerficationEnabled = &isVerificationEnabled
	}

	return &AuthHandler{db, domain, apiKey, opts}
}

func connectDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// AutoMigrate the models
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Token{})

	return db
}

func generateEmailVerificationToken() (string, error) {
	tokenBytes := make([]byte, 64)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func (h *AuthHandler) sendVerificationEmail(receipt string, token string) error {
	USER := "test-pm"
	from := mail.NewEmail(USER, "monedero-luna-saturno@em2826.deviloza.com.mx")
	subject := "Please verify your email"
	to := mail.NewEmail(USER, receipt)
	URI := fmt.Sprintf("%sverify?token=%s", h.Domain, token)
	htmlContent := fmt.
		Sprintf("<p>Click on the <a href=%s>link</a> to verify your email.</p><br>", URI)
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient(h.SendGridAPIKey)
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	log.Println(response.StatusCode)
	if response.StatusCode != 202 {
		log.Println(response.Body)
		// log.Println(response)
		return errors.New(response.Body)
	}
	return nil
}

func (h *AuthHandler) SignUpUser(c echo.Context) error {
	var existingUser User
	user := new(User)

	if err := c.Bind(user); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	err := h.DBConn.Where("email = ?", user.Email).First(&existingUser).Error
	if err == nil {
		log.Println(err)
		log.Println(existingUser)
		msg := fmt.Sprintf("Email %s is already taken", user.Email)
		return errors.New(msg)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to hash password",
		)
	}
	user.Password = string(hashedPassword)

	h.DBConn.Create(&user)
	u := UserDTO{
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
		IsVerified: false,
	}

	// Generate and send a verification email
	encodedToken, err := generateEmailVerificationToken()
	expiration := time.Now().Add(time.Hour)
	if err != nil {
		log.Println(err)
		// Delete previously created user because the client won't be able to verify this account
		h.DBConn.Unscoped().Delete(&user)
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to send the verification mail",
		)
	}

	if *h.Options.IsEmailVerficationEnabled {
		err = h.sendVerificationEmail(u.Email, encodedToken)
		if err != nil {
			log.Println(err)
			// Delete previously created user because the client won't be able to verify this account
			h.DBConn.Unscoped().Delete(&user)
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				"Failed to send the verification mail",
			)
		}
	}

	token := &Token{
		EncodedToken:   encodedToken,
		ExpirationDate: &expiration,
		UserID:         user.ID,
	}

	h.DBConn.Create(token)
	h.DBConn.Model(&user).Update("token_id", token.ID)

	if *h.Options.IsEmailVerficationEnabled {
		return c.JSON(http.StatusOK, u)
	} else {
		jwt, err := GenerateJWT(user.ID, user.Email)
		if err != nil {
			log.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
		}
		return c.JSON(http.StatusOK, JWTResponse{jwt})
	}
}

func (h *AuthHandler) ValidateEmail(c echo.Context) error {
	var verificationToken Token
	var user User

	reqToken := c.QueryParam("token")
	err := h.DBConn.Where(
		"encoded_token = ?",
		reqToken,
	).Take(&verificationToken).Error
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	// Validate expiration date
	if time.Now().After(*verificationToken.ExpirationDate) {
		return echo.NewHTTPError(http.StatusBadRequest, "Token expired")
	}

	// If token is valid then update user row
	err = h.DBConn.
		Where("id = ? AND is_verified = ?", verificationToken.UserID, false).
		Take(&user).
		Error
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	h.DBConn.Model(&user).Update("is_verified", true)
	h.DBConn.Delete(&verificationToken)

	jwt, err := GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	return c.JSON(http.StatusOK, JWTResponse{jwt})
}

func (h *AuthHandler) LoginUser(c echo.Context) error {
	var userDB User
	user := new(User)

	if err := c.Bind(user); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	err := h.DBConn.Where("email = ?", user.Email).First(&userDB).Error
	if err != nil {
		log.Println(err)
		msg := fmt.Sprintf("There are no users with email %s", user.Email)
		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(userDB.Password),
		[]byte(user.Password),
	)
	if err != nil {
		log.Println(err)
		msg := "Incorrect password"
		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}
	jwt, err := GenerateJWT(userDB.ID, user.Email)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	return c.JSON(http.StatusOK, JWTResponse{jwt})
}
