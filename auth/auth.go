package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/spf13/viper"
)

type AccessToken struct {
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn string `json:"expires_in"`
}

type cognitoClaims struct {
	CognitoGroups []string `json:"cognito:groups"`
	CognitoScopes string   `json:"scope"`
	jwt.StandardClaims
}

func VerifyAccessToken(tokenString string) error {
	publicKeySet, err := jwk.Fetch(context.Background(), viper.GetString("COGNITO_KEYS_URL"))
	if err != nil {
		return err
	}
	tokenClaims, err := jwt.ParseWithClaims(
		tokenString,
		&cognitoClaims{}, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodRSA)

			if !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, errors.New("kid header not found")
			}

			_, ok = token.Claims.(*cognitoClaims)
			if !ok {
				return nil, errors.New("There is problem to get claims")
			}

			keys, ok := publicKeySet.LookupKeyID(kid)
			if !ok {
				return nil, fmt.Errorf("key %v not found", kid)
			}

			var tokenKey interface{}
			if err := keys.Raw(&tokenKey); err != nil {
				return nil, errors.New("failed to create token key")
			}

			return tokenKey, nil
		},
	)

	if err != nil {
		return err
	}

	if !tokenClaims.Valid {
		return errors.New("!tokenClaims.Valid")
	}

	if tokenClaims.Method.Alg() != "RS256" {
		return errors.New("tokenClaims.Method.Alg()")
	}

	if !isCognitoGroupsValid(tokenClaims.Claims.(*cognitoClaims).CognitoGroups) &&
		!isScopeValid(tokenClaims.Claims.(*cognitoClaims).CognitoScopes) {
		return errors.New("!isCognitoGroupsValid && !isScopeValid \n" +
			tokenClaims.Claims.(*cognitoClaims).CognitoScopes + "\n",
		)
	}

	if !isIssuerValid(tokenClaims.Claims.(*cognitoClaims).Issuer) {
		return errors.New("!isIssuerValid")
	}
	return nil
}

func isScopeValid(scope string) bool {
	validScope := strings.Contains(scope, viper.GetString("COGNITO_SCOPE"))
	return validScope
}

func isCognitoGroupsValid(groups []string) bool {
	validGroup := SliceContains(groups, viper.GetString("COGNITO_USER_GROUP"))
	return validGroup
}

func isIssuerValid(iss string) bool {
	return iss == viper.GetString("COGNITO_ISSUER")
}

func SliceContains(arr []string, target string) bool {
	for _, el := range arr {
		if el == target {
			return true
		}
	}

	return false
}
