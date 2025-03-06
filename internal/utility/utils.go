package utility

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"strconv"

	"encoding/base64"
	"regexp"
	"strings"

	"github.com/gboliknow/bildwerk/internal/models"
	"github.com/gin-gonic/gin"
)

func WriteJSON(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := models.Response{
		StatusCode:   statusCode,
		IsSuccessful: true,
		Message:      message,
		Data:         data,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func RespondWithError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"isSuccessful": false,
		"message":      message,
	})
}

func GetTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	headerParts := strings.Split(authHeader, " ")
	errTokenMissing := errors.New("missing or invalid token")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errTokenMissing
	}
	tokenAuth := headerParts[1]
	if tokenAuth != "" {
		return tokenAuth, nil
	}
	return "", errTokenMissing
}

func GenerateResetToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	return fmt.Sprintf("%x", b)
}

func GetQueryFloat(c *gin.Context, key string, defaultValue float64) float64 {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetQueryInt retrieves an integer query parameter from the request, returning a default value if missing or invalid.
func GetQueryInt(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func IsValidBase64Image(base64Str string) bool {
	if len(base64Str) == 0 {
		return false
	}
	pattern := regexp.MustCompile(`^data:image\/([a-zA-Z0-9]+);base64,(.+)$`)
	matches := pattern.FindStringSubmatch(base64Str)

	var actualBase64 string
	if len(matches) == 3 {

		actualBase64 = matches[2]
	} else {

		actualBase64 = base64Str
	}
	actualBase64 = strings.TrimSpace(actualBase64)

	_, err := base64.StdEncoding.DecodeString(actualBase64)
	if err != nil {
		return false
	}
	if len(actualBase64) < 100 {
		return false
	}

	return true
}
