package service

import (
 "fmt"
 "github.com/golang-jwt/jwt/v5"
 "time"
)

var secretKey = []byte("secret-key")

func CreateToken(username string, permission_kind string, permission_name string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
        "username": username, 	// subject scope id
		"permission_kind": permission_kind,
		"permission_name": permission_name,
        "exp": time.Now().Add(time.Hour * 24).Unix(), 
        })

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
    return "", err
    }

 return tokenString, nil
}

func ParseToken(token_string string) {
	token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		// Handle error, e.g., invalid token
		fmt.Println("Error parsing token:", err)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if permission_kind, ok := claims["permission_kind"].(string); ok {
			fmt.Println("Custom Claim 1:", permission_kind)
		} else {
			fmt.Println("Custom Claim 1 is not a string or does not exist.")
		}
	} else {
		fmt.Println("Invalid claims type.")
	}
}