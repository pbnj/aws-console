package console

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/config"
)

const (
	consoleURL = "https://console.aws.amazon.com/"
	signinURL  = "https://signin.aws.amazon.com/federation"
)

// Session represents an AWS Federation session
type Session struct {
	SessionID    string `json:"sessionId"`
	SessionKey   string `json:"sessionKey"`
	SessionToken string `json:"sessionToken"`
}

// Token represents AWS SigninToken object
type Token struct {
	SigninToken string `json:"SigninToken"`
}

// URL returns the AWS Console URL
func URL(userName, awsProfile string) (string, error) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(awsProfile))
	creds, err := cfg.Credentials.Retrieve(ctx)

	loginSession := &Session{
		SessionID:    creds.AccessKeyID,
		SessionKey:   creds.SecretAccessKey,
		SessionToken: creds.SessionToken,
	}

	loginSessionJSON, err := json.Marshal(loginSession)
	if err != nil {
		return "", fmt.Errorf("could not marshal JSON: %v", err)
	}

	signinTokenURL := signinURL +
		"?Action=getSigninToken" +
		"&SessionType=json" +
		"&Session=" +
		url.QueryEscape(string(loginSessionJSON))

	resp, err := http.Get(signinTokenURL)
	if err != nil {
		return "", fmt.Errorf("could not get sign-in token URL: %v", err)
	}
	defer resp.Body.Close()

	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var token Token
	err = json.Unmarshal(respJSON, &token)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal JSON token: %v", err)
	}

	awsConsoleURL := signinURL +
		"?Action=login" +
		"&SigninToken=" + url.QueryEscape(token.SigninToken) +
		"&Issuer=" + userName +
		"&Destination=" + url.QueryEscape(consoleURL)

	return awsConsoleURL, nil
}
