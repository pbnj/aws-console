package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	console "github.com/pbnj/aws-console"
	open "github.com/pbnj/go-open"
)

// Version is the version of this Go binary
const Version = "0.1.0"

func main() {
	var (
		help       = flag.Bool("h", false, "Print Help")
		version    = flag.Bool("v", false, "Print Version")
		awsProfile string

		Usage = func() {
			fmt.Printf("Usage of %s:\n", os.Args[0])
			flag.PrintDefaults()
			os.Exit(1)
		}
		currentUser *user.User
		userName    string
	)

	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err.Error())
	}

	userName = currentUser.Username
	flag.StringVar(&awsProfile, "p", "", "AWS Profile")

	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if *help {
		Usage()
		os.Exit(0)
	}

	if awsProfile == "" {
		log.Println("-p is required")
		Usage()
	}

	awsConsoleURL, err := console.URL(userName, awsProfile)
	if err != nil {
		log.Fatal(err)
	}
	open.Open(awsConsoleURL)
}
