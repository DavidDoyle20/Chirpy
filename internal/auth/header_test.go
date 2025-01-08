package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBearerToken(t *testing.T) {
	// Empty headers
	token, err := GetBearerToken(http.Header{})
	assert.Error(t, err, "Should error when there is no header")
	assert.Equal(t, token, "", "Token should be blank if there is an error")

	// One header that isnt the correct one
	headers := http.Header{}
	headers.Add("Test", "test")
	token, err = GetBearerToken(headers)
	assert.Error(t, err, "Should error when there is no bearer key in the header")
	assert.Equal(t, token, "", "Token should be blank if there is no header")

	// Header contains Bearer key
	headers.Add("Bearer", "jwt")
	token, err = GetBearerToken(headers)
	assert.NoError(t, err, "Should not error when there is a bearer key in the header")
	assert.Equal(t, token, "jwt", "The correct jwt should be returned")
}
