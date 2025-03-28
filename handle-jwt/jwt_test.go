package handlejwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	token, err := GenerateJWTToken("1", "admin")
	assert.NoError(t, err)

	// Expired를 테스트 하기 위함
	//time.Sleep(2 * time.Minute)

	claims, err := ValidateJWT(token)
	assert.NoError(t, err)

	assert.Equal(t, "1", claims.RegisteredClaims.Subject)
	assert.Equal(t, "myserver.com", claims.RegisteredClaims.Issuer)
	assert.Equal(t, "admin", claims.Role)
}
