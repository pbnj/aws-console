package console

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/fatih/color"
)

const (
	consoleURL = "https://console.aws.amazon.com/"
	signinURL  = "https://signin.aws.amazon.com/federation"
)

// Session represents AWS Console session
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
func URL(userName, awsCredFileFlag, awsProfileFlag string) (string, error) {
	creds := credentials.NewSharedCredentials(awsCredFileFlag, awsProfileFlag)

	sess := session.Must(session.NewSession())
	svc := sts.New(sess, aws.NewConfig().WithCredentials(creds))

	tokenOutput, err := svc.GetFederationToken(
		&sts.GetFederationTokenInput{
			Name: aws.String(userName),
			Policy: aws.String(
				`{
					"Version": "2012-10-17",
					"Statement": [
						{
							"Sid": "Stmt1437414476731",
							"Action": "*",
							"Effect": "Allow",
							"Resource": "*"
						}
					]
				}`,
			),
		})
	if err != nil {
		return "", fmt.Errorf("could not get federation token: %v", err)
	}

	logrus.Debugf("Credentials from AWS Federation: %+v\n", *tokenOutput.Credentials)

	loginSession := &Session{
		SessionID:    *tokenOutput.Credentials.AccessKeyId,
		SessionKey:   *tokenOutput.Credentials.SecretAccessKey,
		SessionToken: *tokenOutput.Credentials.SessionToken,
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

	logrus.Debugf("Signin Token URL: %+v\n", signinTokenURL)

	resp, err := http.Get(signinTokenURL)
	if err != nil {
		return "", fmt.Errorf("could not get sign-in token URL: %v", err)
	}
	defer resp.Body.Close()

	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red(err.Error())
	}

	var token Token
	err = json.Unmarshal(respJSON, &token)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal JSON token: %v", err)
	}

	logrus.Debugf("Token from AWS Federation: %+v\n", token)

	awsConsoleURL := signinURL + "?Action=login" + "&SigninToken=" + url.QueryEscape(token.SigninToken) + "&Issuer=" + userName + "&Destination=" + url.QueryEscape(consoleURL)

	logrus.Debugf("AWS Console URL to be launched: %+v\n", awsConsoleURL)

	return awsConsoleURL, nil
}
