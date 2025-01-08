package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	vals := headers.Values("Authorization")
	tokenString := ""
	for _, v := range vals {
		if strings.Contains(v, "Bearer") {
			tokenString = strings.TrimSpace(strings.Replace(v, "Bearer ", "", 1))
		}
	}
	if tokenString == "" {
		return "", fmt.Errorf("bearer header doesnt exist")
	}
	return tokenString, nil
}