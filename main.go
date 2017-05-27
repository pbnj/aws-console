package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/fatih/color"
	flag "github.com/ogier/pflag"
	open "github.com/petermbenjamin/go-open"
)

var (
	debugFlag       = flag.BoolP("debug", "d", false, "Debug")
	helpFlag        = flag.BoolP("help", "h", false, "Print Help")
	awsCredFileFlag string
	awsProfileFlag  string

	currentUser *user.User
	userName    string
	homeDir     string
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

func setup() {
	currentUser, err := user.Current()
	if err != nil {
		color.Red(err.Error())
	}

	userName = currentUser.Username
	homeDir = currentUser.HomeDir
}
func init() {
	setup()
	flag.StringVarP(&awsCredFileFlag, "credentials", "c", filepath.Join(homeDir, ".aws", "credentials"), "Path to AWS credentials file")
	flag.StringVarP(&awsProfileFlag, "profile", "p", "default", "AWS Profile")
}

func main() {
	flag.Parse()

	if *debugFlag {
		log.SetLevel(log.DebugLevel)
	}

	if *helpFlag {
		fmt.Println("Usage: aws-console <options>")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	creds := credentials.NewSharedCredentials(awsCredFileFlag, awsProfileFlag)

	sess := session.Must(session.NewSession())
	svc := sts.New(sess, aws.NewConfig().WithCredentials(creds))

	tokenOutput, err := svc.GetFederationToken(&sts.GetFederationTokenInput{Name: aws.String(userName), Policy: aws.String("{\"Version\": \"2012-10-17\", \"Statement\": [{\"Sid\": \"Stmt1437414476731\", \"Action\": \"*\",\"Effect\": \"Allow\", \"Resource\": \"*\" }]}")})
	if err != nil {
		color.Red(err.Error())
	}

	log.Debugf("Credentials from AWS Federation: %+v\n", *tokenOutput.Credentials)

	loginSession := &Session{
		SessionID:    *tokenOutput.Credentials.AccessKeyId,
		SessionKey:   *tokenOutput.Credentials.SecretAccessKey,
		SessionToken: *tokenOutput.Credentials.SessionToken,
	}

	loginSessionJSON, err := json.Marshal(loginSession)
	if err != nil {
		color.Red(err.Error())
	}

	signinTokenURL := signinURL + "?Action=getSigninToken" + "&SessionType=json" + "&Session=" + url.QueryEscape(string(loginSessionJSON))

	log.Debugf("Signin Token URL: %+v\n", signinTokenURL)

	resp, err := http.Get(signinTokenURL)
	if err != nil {
		color.Red(err.Error())
	}
	defer resp.Body.Close()

	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Red(err.Error())
	}

	var token Token
	err = json.Unmarshal(respJSON, &token)
	if err != nil {
		color.Red(err.Error())
	}

	log.Debugf("Token from AWS Federation: %+v\n", token)

	awsConsoleURL := signinURL + "?Action=login" + "&SigninToken=" + url.QueryEscape(token.SigninToken) + "&Issuer=" + userName + "&Destination=" + url.QueryEscape(consoleURL)

	log.Debugf("AWS Console URL to be launched: %+v\n", awsConsoleURL)

	open.Open([]string{awsConsoleURL})
}
