package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── HashPassword ────────────────────────────────────────────────────────────

func TestHashPassword_GeneraHashDiferenteAlOriginal(t *testing.T) {
	hash, err := HashPassword("miPassword123")
	require.NoError(t, err)
	assert.NotEqual(t, "miPassword123", hash)
}

func TestHashPassword_MismaContraseñaGeneraHashesDistintos(t *testing.T) {
	hash1, _ := HashPassword("miPassword123")
	hash2, _ := HashPassword("miPassword123")
	assert.NotEqual(t, hash1, hash2, "bcrypt debe generar salt distinto cada vez")
}

func TestHashPassword_ContraseñaVaciaRetornaError(t *testing.T) {
	_, err := HashPassword("")
	assert.Error(t, err)
}

// ─── CheckPassword ───────────────────────────────────────────────────────────

func TestCheckPassword_ContraseñaCorrectaRetornaTrue(t *testing.T) {
	hash, _ := HashPassword("correcta123")
	ok := CheckPassword("correcta123", hash)
	assert.True(t, ok)
}

func TestCheckPassword_ContraseñaIncorrectaRetornaFalse(t *testing.T) {
	hash, _ := HashPassword("correcta123")
	ok := CheckPassword("incorrecta456", hash)
	assert.False(t, ok)
}

func TestCheckPassword_HashVacioRetornaFalse(t *testing.T) {
	ok := CheckPassword("cualquier", "")
	assert.False(t, ok)
}

// ─── GenerateJWT ─────────────────────────────────────────────────────────────

func TestGenerateJWT_TokenContieneUserID(t *testing.T) {
	token, err := GenerateJWT(42, "user@test.com", "admin", "test-secret")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claims, err := ValidateJWT(token, "test-secret")
	require.NoError(t, err)
	assert.Equal(t, 42, claims.UserID)
}

func TestGenerateJWT_TokenContieneRol(t *testing.T) {
	token, _ := GenerateJWT(1, "user@test.com", "empleado", "test-secret")
	claims, err := ValidateJWT(token, "test-secret")
	require.NoError(t, err)
	assert.Equal(t, "empleado", claims.Rol)
}

func TestGenerateJWT_TokenContieneEmail(t *testing.T) {
	token, _ := GenerateJWT(1, "user@test.com", "admin", "test-secret")
	claims, err := ValidateJWT(token, "test-secret")
	require.NoError(t, err)
	assert.Equal(t, "user@test.com", claims.Email)
}

func TestGenerateJWT_TokenExpiraEn8Horas(t *testing.T) {
	token, _ := GenerateJWT(1, "user@test.com", "admin", "test-secret")
	claims, err := ValidateJWT(token, "test-secret")
	require.NoError(t, err)

	expiracion := time.Unix(claims.ExpiresAt.Unix(), 0)
	diferencia := expiracion.Sub(time.Now())

	// Debe expirar entre 7h55m y 8h05m desde ahora
	assert.Greater(t, diferencia, 7*time.Hour+55*time.Minute)
	assert.Less(t, diferencia, 8*time.Hour+5*time.Minute)
}

// ─── ValidateJWT ─────────────────────────────────────────────────────────────

func TestValidateJWT_TokenValidoDevuelveClaimsCorrectos(t *testing.T) {
	token, _ := GenerateJWT(99, "admin@res.com", "admin", "secret123")
	claims, err := ValidateJWT(token, "secret123")

	require.NoError(t, err)
	assert.Equal(t, 99, claims.UserID)
	assert.Equal(t, "admin@res.com", claims.Email)
	assert.Equal(t, "admin", claims.Rol)
}

func TestValidateJWT_TokenMalformadoDevuelveError(t *testing.T) {
	_, err := ValidateJWT("esto.no.es.un.jwt.valido", "secret")
	assert.Error(t, err)
}

func TestValidateJWT_TokenConFirmaFalsaDevuelveError(t *testing.T) {
	token, _ := GenerateJWT(1, "user@test.com", "admin", "secret-correcto")
	_, err := ValidateJWT(token, "secret-falso")
	assert.Error(t, err)
}

func TestValidateJWT_TokenVacioDevuelveError(t *testing.T) {
	_, err := ValidateJWT("", "secret")
	assert.Error(t, err)
}
