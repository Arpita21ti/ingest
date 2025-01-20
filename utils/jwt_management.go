package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecretKey []byte

// Initialize the JWT secret key.
func init() {
	var err error
	jwtSecretKey, err = loadJWTKey()
	if err != nil {
		panic(fmt.Sprintf("Failed to load JWT key: %v", err))
	}
}

// Function to generate a random key.
func generateRandomKey() ([]byte, error) {
	// Generates 32 random bytes (256-bit key for strong security)
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random key: %v", err)
	}
	// Optionally, encode the key as a base64 string for storage
	// encodedKey := base64.StdEncoding.EncodeToString(key)
	// return []byte(encodedKey), nil

	return key, nil
}

// Function to load the JWT secret key from the environment (if available)
// or generate a new one if it's not set.
func loadJWTKey() ([]byte, error) {
	// First, check if the key is stored in an environment variable
	keyFromEnv := os.Getenv("JWT_SECRET_KEY")
	if keyFromEnv != "" {
		return []byte(keyFromEnv), nil
	}
	// If no key is found in the environment, generate a random key
	key, err := generateRandomKey()
	if err != nil {
		return nil, err
	}
	// [Optional]: stored the generated key in the environment for the current session.
	os.Setenv("JWT_SECRET_KEY", base64.StdEncoding.EncodeToString(key))
	return key, nil
}

// GenerateToken creates a JWT with user role or privileges
func GenerateToken(username, role string, expiryHours float32) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	expirationTime := time.Now().Add(time.Duration(expiryHours) * time.Hour).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["username"] = username
	claims["role"] = role // Add role to token claims
	claims["exp"] = expirationTime

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("can't generate JWT token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates the JWT and checks privileges
func ValidateToken(tokenString string) (jwt.MapClaims, error) {

	// Parse and validate the token.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Ensure the token is valid.
	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// Extract claims.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse token claims")
	}

	// Check expiration claim.
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return nil, fmt.Errorf("token expired")
		}
	} else {
		return nil, fmt.Errorf("missing or invalid expiration claim")
	}

	return claims, nil
}

// Middleware to validate privileges
func HasPrivilege(claims jwt.MapClaims, requiredRole string) bool {
	role, ok := claims["role"].(string)
	if !ok {
		return false // No role present
	}

	// Define role hierarchy
	roleHierarchy := map[string]int{
		"common": 1, // Basic users
		"admin":  2, // Admin users
		"master": 3, // Master users
	}

	// Check if the user's role has sufficient privilege
	return roleHierarchy[role] >= roleHierarchy[requiredRole]
}

// // Generate tokens for different roles
// commonToken, err := GenerateToken("student_user", "common", 24)
// if err != nil {
// 	fmt.Println("Error generating common user token:", err)
// }

// adminToken, err := GenerateToken("admin_user", "admin", 24)
// if err != nil {
// 	fmt.Println("Error generating admin token:", err)
// }

// masterToken, err := GenerateToken("owner_user", "master", 24)
// if err != nil {
// 	fmt.Println("Error generating master token:", err)
// }
