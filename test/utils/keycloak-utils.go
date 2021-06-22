package utils

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/dgrijalva/jwt-go"
)

func GetToken() string {
	keycloakURL := GetConfig("../test-cicd.toml").Env.KeycloakURL
	username := GetConfig("../test-cicd.toml").Env.Username
	password := GetConfig("../test-cicd.toml").Env.Password
	clientID := GetConfig("../test-cicd.toml").Env.ClientID
	getTokenURL := fmt.Sprintf("%v/realms/develop/protocol/openid-connect/token", keycloakURL)
	resp, err := http.PostForm(getTokenURL,
		url.Values{"username": {username}, "password": {password}, "client_id": {clientID}, "grant_type": {"password"}, "scope": {"openid"}})
	if err != nil {
		fmt.Println(err)
		panic("post request to keycloak failed.")
	}
	defer resp.Body.Close()
	returnMap, err := ParseResponse(resp)
	if err != nil {
		fmt.Println(err)
		panic("cannot parse response from keycloak.")
	}
	if str, ok := returnMap["id_token"].(string); ok {
		return str
	}
	return ""
}

func GetUserID(tokenString string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		subString := fmt.Sprintf("%v", claims["sub"])
		return subString, nil
	} else {
		return "", err
	}
	// return "test", err
	// _, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	claims := token.Claims.(jwt.MapClaims)
	// 	return claims["sub"], nil
	// 	fmt.Println(claims["sub"])
	// 	//Make sure that the token method conform to "SigningMethodRS256"
	// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 		claims := token.Claims.(jwt.MapClaims)
	// 		// data := claims["data"].(map[string]interface{})
	// 		// userID := data["sub"].(string)
	// 		fmt.Println("aa")
	// 		fmt.Println(claims["sub"])
	// 		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	// 	}
	// 	return []byte(os.Getenv("ACCESS_SECRET")), nil
	// })
	// return "", err
}
