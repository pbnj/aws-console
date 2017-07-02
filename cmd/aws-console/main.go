package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	flag "github.com/ogier/pflag"
	console "github.com/petermbenjamin/aws-console"
	open "github.com/petermbenjamin/go-open"
)

// Version is the version of this Go binary
const Version = "0.0.1"

var (
	debugFlag       = flag.BoolP("debug", "d", false, "Debug")
	helpFlag        = flag.BoolP("help", "h", false, "Print Help")
	versionFlag     = flag.BoolP("version", "v", false, "Print Version")
	awsCredFileFlag string
	awsProfileFlag  string

	currentUser *user.User
	userName    string
	homeDir     string
)

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

	if *versionFlag {
		fmt.Println(Version)
		os.Exit(0)
	}

	if *helpFlag {
		fmt.Println("Usage: aws-console <options>")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	awsConsoleURL, err := console.URL(userName, awsCredFileFlag, awsProfileFlag)
	if err != nil {
		color.Red(err.Error())
	}
	open.Open(awsConsoleURL)
}
