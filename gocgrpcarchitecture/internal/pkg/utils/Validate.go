package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Escape(value string) string {
	return html.EscapeString(strings.TrimSpace(IsTitle(value)))
}
func IsTitle(value string) string {
	return cases.Title(language.English, cases.Compact).String(value)
}
func IsToLower(value string) string {
	return cases.Lower(language.English, cases.Compact).String(value)
}

func IsMatch(value string) (string, error) {
	//US := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	//US := regexp.MustCompile(`^\(?(\d{3})\)?[-.\s]?(\d{3})[-.\s]?(\d{4})$`)
	IntelAndLocal := regexp.MustCompile(`^\+?[0-9]\d{0,3}[-.\s]?\d{1,14}$`)
	if !IntelAndLocal.MatchString(value) {
		return "", errors.New("numeric format must be valid")
	}
	return value, nil
}

func RegexValidate(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) < 2 || len(value) > 45 {
		return "", errors.New("field must be between 2 and 45 characters")
	}
	re := regexp.MustCompile(`^[a-zA-Z]+$`) // match without space
	if !re.MatchString(value) {
		return "", errors.New("character must be valid")
	}
	return value, nil
}

func RegexValidatePhone(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) < 2 || len(value) > 45 {
		return "", errors.New("field must be between 2 and 45 characters")
	}
	//re := regexp.MustCompile(`^[0-9-]+$`) // match without space
	//re := regexp.MustCompile(`^\d{3}\d{3}\d{4}$`) // match without space
	re := regexp.MustCompile(`^\+[1-9][1-9][0-9]{1,14}$`) // match without space
	if !re.MatchString(value) {
		return "", errors.New("character must be valid")
	}
	return value, nil
}
