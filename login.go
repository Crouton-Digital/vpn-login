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

type LogFile struct {
	File *os.File
}

func main() {
	log := AccessLog()
	defer log.File.Close()

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
		log.Write(err.Error() + "\n")
		os.WriteFile(os.Getenv("auth_control_file"), []byte("0"), 0644)
		os.Exit(0)
	}

	app.Token = *authResp.AuthenticationResult.AccessToken
	log.Write("Success:" + app.Token + "\n")

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

func AccessLog() *LogFile {
	f, err := os.OpenFile(os.Getenv("LOG_FILE"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return &LogFile{File: f}
}

func (log *LogFile) Write(t string) {
	if os.Getenv("LOG_ENABLED") == "1" {
		_, err := log.File.WriteString(t)
		panic(err)
	}
}
