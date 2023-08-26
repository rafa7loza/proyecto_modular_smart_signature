package web

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"deviloza.com.mx/auth"
)

type WebHandler struct {
	DBConn         *gorm.DB
	AvailableHosts *Hosts
}

func NewWebHandler(dsn string, hosts *Hosts) *WebHandler {
	db := connectDB(dsn)
	return &WebHandler{db, hosts}
}

func connectDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Document{}, &ProcessedDocument{})
	return db
}

func (h *WebHandler) Profile(c echo.Context) error {
	var user auth.User
	reqToken := c.Get("user").(*jwt.Token)
	log.Println(reqToken)
	claims := reqToken.Claims.(*auth.JWTCustomClaims)
	log.Println(claims)

	h.DBConn.Take(&user)
	log.Println(user)

	return c.JSON(http.StatusOK, user)
}

func (h *WebHandler) UploadFile(c echo.Context) error {
	file, err := c.FormFile("files")
	reqToken := c.Get("user").(*jwt.Token)

	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	fileName := file.Filename
	ext := filepath.Ext(fileName)
	if ext == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	src, err := file.Open()
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}
	defer src.Close()

	content, err := io.ReadAll(src)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	// Get the claims from the token
	claims := reqToken.Claims.(*auth.JWTCustomClaims)
	doc := &Document{
		DocumentContent: content,
		FileName:        fileName[:len(fileName)-len(ext)],
		Extension:       ext,
		UserId:          claims.UID,
	}
	h.DBConn.Create(doc)

	if err := h.startProcess(doc.ID); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal error")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"Message": "File uploaded successfully",
	})
}

func (h *WebHandler) startProcess(docId uint) error {
	curHost, err := h.AvailableHosts.GetNext()
	if err != nil {
		return err
	}

	log.Println(curHost)
	endpointURL := fmt.Sprintf(
		"http://%s:%s/detect/%d",
		curHost.Address,
		curHost.Port,
		docId,
	)

	log.Println(endpointURL)
	req, err := http.NewRequest(http.MethodPost, endpointURL, nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	log.Println(res)
	return nil
}

func (h *WebHandler) GetDocument(c echo.Context) error {
	docId := c.Param("docId")
	var document Document

	err := h.DBConn.Where("id = ?", docId).First(&document).Error
	if err != nil {
		log.Println(err)
		msg := fmt.Sprintf("Document with id %d not found", docId)
		return echo.NewHTTPError(http.StatusBadRequest, msg)
	}

	// encodedContent += "data:image/jpeg;base64,"
	encodedContent := base64.StdEncoding.EncodeToString(document.DocumentContent)
	res := DocumentRes{
		document.FileName,
		document.Extension,
		"base64",
		encodedContent,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *WebHandler) GetUserDocuments(c echo.Context) error {
	var docs []Document
	var docsRes DocumentsRes

	reqToken := c.Get("user").(*jwt.Token)
	claims := reqToken.Claims.(*auth.JWTCustomClaims)

	err := h.DBConn.Where("user_id = ?", claims.UID).Find(&docs).Error
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request")
	}

	for _, doc := range docs {
		log.Println(doc.FileName, doc.Extension)
		encodedContent := base64.StdEncoding.EncodeToString(doc.DocumentContent)
		d := DocumentRes {
			doc.FileName,
			doc.Extension,
			"base64",
			encodedContent,
		}
		docsRes.Documents = append(docsRes.Documents, d)
	}

	return c.JSON(http.StatusOK, docsRes)
}
