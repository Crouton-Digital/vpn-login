// Package: main
// Language: go
// Path: login.go
// Author: Rolands Kalpins
// Date: 2022-08-09

package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

type App struct {
	CognitoClient   *cognito.CognitoIdentityProvider
	UserPoolID      string
	AppClientID     string
	AppClientSecret string
	Token           string
}

func main() {

	app := NewCognito()

	authResp, err := app.CognitoClient.InitiateAuth(
		&cognito.InitiateAuthInput{
			AuthFlow: aws.String("USER_PASSWORD_AUTH"),
			AuthParameters: map[string]*string{
				"USERNAME":    aws.String(os.Getenv("username")),
				"PASSWORD":    aws.String(os.Getenv("password")),
				"SECRET_HASH": aws.String(computeSecretHash(app.AppClientSecret, os.Getenv("username"), app.AppClientID)),
			},
			ClientId: aws.String(app.AppClientID),
		},
	)
	if err != nil {
		os.WriteFile(os.Getenv("auth_control_file"), []byte("0"), 0644)
		os.Exit(0)
	}

	app.Token = *authResp.AuthenticationResult.AccessToken

	os.WriteFile(os.Getenv("auth_control_file"), []byte("1"), 0644)
	os.Exit(0)
}

func computeSecretHash(clientSecret string, username string, clientId string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func NewCognito() *App {
	conf := &aws.Config{Region: aws.String(os.Getenv("COGNITO_REGION"))}
	mySession := session.Must(session.NewSession(conf))
	return &App{
		CognitoClient:   cognito.New(mySession),
		UserPoolID:      os.Getenv("COGNITO_USER_POOL_ID"),
		AppClientID:     os.Getenv("COGNITO_APP_CLIENT_ID"),
		AppClientSecret: os.Getenv("COGNITO_APP_CLIENT_SECRET"),
	}
}

// f, err := os.OpenFile("/tmp/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// if err != nil {
// 	log.Fatal(err)
// }
// defer f.Close()

// var arguments string
// for _, arg := range os.Args[:] {
// 	arguments += arg + " "
// }

// arguments += " *** With ENV ***"
// arguments += os.Getenv("username")
// arguments += " " + os.Getenv("password")

// arguments += "\n"

// if _, err := f.Write([]byte(arguments)); err != nil {
// 	f.Close() // ignore error; Write error takes precedence
// 	log.Fatal(err)
// }
