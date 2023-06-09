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
	"time"

	"bitbucket.org/am-bitbucket/vpn-login/auth"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/dgrijalva/jwt-go/v4"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
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

type cognitoClaims struct {
	CognitoGroups []string `json:"cognito:groups"`
	CognitoScopes string   `json:"scope"`
	jwt.StandardClaims
}

func main() {
	ReadConfigFile()
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
		log.Write(os.Getenv("username") + ":" + err.Error())
		os.WriteFile(os.Getenv("auth_control_file"), []byte("0"), 0644)
		os.Exit(0)
	}

	app.Token = *authResp.AuthenticationResult.AccessToken

	if err := auth.VerifyAccessToken(app.Token); err != nil {
		log.Write(os.Getenv("username") + ":" + err.Error())
		os.WriteFile(os.Getenv("auth_control_file"), []byte("0"), 0644)
		os.Exit(0)
	}

	log.Write("Access granted " + os.Getenv("username"))

	os.WriteFile(os.Getenv("auth_control_file"), []byte("1"), 0644)
	os.Exit(0)
}

func computeSecretHash(clientSecret string, username string, clientId string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func NewCognito() *App {
	conf := &aws.Config{Region: aws.String(viper.GetString("COGNITO_REGION"))}
	mySession := session.Must(session.NewSession(conf))
	return &App{
		CognitoClient:   cognito.New(mySession),
		UserPoolID:      viper.GetString("COGNITO_USER_POOL_ID"),
		AppClientID:     viper.GetString("COGNITO_APP_CLIENT_ID"),
		AppClientSecret: viper.GetString("COGNITO_APP_CLIENT_SECRET"),
	}
}

func AccessLog() *LogFile {
	f, err := os.OpenFile(viper.GetString("LOG_FILE"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// f, err := os.OpenFile("/tmp/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return &LogFile{File: f}
}

func (log *LogFile) Write(t string) {
	if viper.GetString("LOG_ENABLED") == "1" {
		_, err := log.File.WriteString(time.Now().String() + ": " + t + "\n")
		if err != nil {
			panic(err)
		}
	}
}

func ReadConfigFile() {
	viper.SetConfigName("settings.yml")
	viper.AddConfigPath("./")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	viper.AutomaticEnv()
}
