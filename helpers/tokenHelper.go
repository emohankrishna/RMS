package helpers

import (
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go/v4"
	"golang.org/x/crypto/bcrypt"
)

func GetHash(str []byte) string {
	hashDigest, err := bcrypt.GenerateFromPassword(str, bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
	}
	return string(hashDigest)
}

func VerifyPassword(hashedPassword, password []byte) bool {
	fmt.Println(string(hashedPassword), string(password))
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	return err == nil
}

type UserClaims struct {
	jwt.StandardClaims
	User_id string
}

func CreateToken(c *UserClaims) (string, error) {
	var jwtKey []byte = []byte(GetEnvString("SECRET_KEY"))
	var token *jwt.Token = jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// claims := token.Claims.(jwt.MapClaims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("error in createToken")
	}
	return tokenString, nil
}

func ParseToken(signedToken string) (*UserClaims, error) {
	var jwtKey []byte = []byte(GetEnvString("SECRET_KEY"))
	token, err := jwt.ParseWithClaims(signedToken, &UserClaims{
		User_id: "bindu@gamil.com",
	}, func(unverifiedToken *jwt.Token) (interface{}, error) {
		fmt.Println("unverifiedToken : ", unverifiedToken)
		if unverifiedToken.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("invalid signing algorithm")
		}
		return jwtKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error in parseToken while parsing token %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("error in parseToken, token is not valid")
	}
	return token.Claims.(*UserClaims), nil
}
