// Copyright 2019 go-azuredevops AUTHORS. All rights reserved.
//
// This example uses the go-azuredevops library to handle a webhook
// from Azure Devops.  Use a tool like ngrok to test locally.
// This code:
// 0. Prompts the user for input
// 1. Creates an API client using basic auth and personal access token
// 2. Sets up a web server on port 9000
// 3. Listens for azuredevops.PushEvent webhooks
// 4. Upon receiving a webhook, validates Basic user/pass
// 5. Queues the supplied build ID number if webhook is valid
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"

	"github.com/mcdafydd/go-azuredevops/azuredevops"
	"golang.org/x/crypto/ssh/terminal"
)

type myClient struct {
	c            *azuredevops.Client
	token        string
	org          string
	project      string
	sourceBranch string
	buildName    string
	buildID      int
	whUser       string
	whPass       string
}

func main() {
	client := new(myClient)
	fmt.Println("The org, project, and branch you supply will be ")
	fmt.Println("given to the Build() call to trigger a new build ")
	fmt.Println("after receiving the webhook.")

	r := bufio.NewReader(os.Stdin)
	fmt.Print("Azure Devops Organization: ")
	client.org, _ = r.ReadString('\n')
	client.org = strings.TrimSuffix(client.org, "\n")

	r = bufio.NewReader(os.Stdin)
	fmt.Print("Azure Devops Project: ")
	client.project, _ = r.ReadString('\n')
	client.project = strings.TrimSuffix(client.project, "\n")

	r = bufio.NewReader(os.Stdin)
	fmt.Print("Azure Devops Branch To Build: ")
	client.sourceBranch, _ = r.ReadString('\n')
	client.sourceBranch = strings.TrimSuffix(client.sourceBranch, "\n")

	fmt.Print("Azure Devops Build ID (must be ID number, not name): ")
	tmp, _ := r.ReadString('\n')
	tmp = strings.TrimSuffix(tmp, "\n")
	buildID, err := strconv.Atoi(tmp)
	client.buildID = buildID
	if err != nil {
		fmt.Printf("Error getting build ID number from user: %+v\n", err)
		return
	}

	r = bufio.NewReader(os.Stdin)
	fmt.Print("Azure Devops Webhook Basic username: ")
	client.whUser, _ = r.ReadString('\n')
	client.whUser = strings.TrimSuffix(client.whUser, "\n")

	r = bufio.NewReader(os.Stdin)
	fmt.Print("Azure Devops Webhook Basic password: ")
	client.whPass, _ = r.ReadString('\n')
	client.whPass = strings.TrimSuffix(client.whPass, "\n")

	fmt.Print("Azure Devops Personal Access Token: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	client.token = string(bytePassword)

	// NOTE: username is an empty string when using the
	// personal access token auth method
	tp := azuredevops.BasicAuthTransport{
		Username: "",
		Password: strings.TrimSpace(client.token),
	}

	client.c, _ = azuredevops.NewClient(tp.Client())

	// Create our logger
	logger := log.New(os.Stdout, "", 0)
	// Register our web server handler.
	http.Handle("/events", eventHandler(client, logger))
	http.ListenAndServe(":9000", nil)
}

func eventHandler(client *myClient, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, err := azuredevops.ValidatePayload(r, []byte(client.whUser), []byte(client.whPass))
		if payload == nil && err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintln(w, "Authentication failed")
			logger.Printf("Webhook auth failed in eventHandler: %s", err)
			return
		} else if err != nil {
			logger.Printf("Err in eventHandler: %s", err)
			return
		}
		event, err := azuredevops.ParseWebHook(payload)
		if err != nil {
			logger.Printf("Err in eventHandler: %s", err)
			return
		}
		if event.PayloadType == azuredevops.PushEvent {
			// If we pass the basic auth, just run the build
			// without doing any other checks on the webhook
			// payload
			buildDef := azuredevops.BuildDefinition{
				ID: Int(1),
			}
			build := azuredevops.Build{
				Definition:   &buildDef,
				SourceBranch: formatRef(client.sourceBranch),
			}
			opts := azuredevops.QueueBuildOptions{}
			r, resp, err := client.c.Builds.Queue(
				context.Background(),
				client.org,
				client.project,
				&build,
				&opts)
			if err != nil {
				logger.Printf("HTTP error: %+v\n", err)
			} else {
				logger.Printf("Successful build trigger response: %+v\n", r)
				logger.Printf("\nHTTP Response: %+v\n", resp)
			}
		}
	})
}

// formatRef helper function for API calls that need a branch reference
// as an input parameter and returns a pointer to a formatted string.
// Examples:
// *ref = "mybranch" => *ref = "refs/heads/mybranch"
// *ref = "refs/heads/abranch" => *ref = "refs/heads/abranch"
func formatRef(ref string) *string {
	matched, err := path.Match("refs/heads/*", ref)
	if matched && err == nil {
		return &ref
	} else if err != nil {
		return nil
	} else {
		ret := fmt.Sprintf("refs/heads/%s", ref)
		return &ret
	}
}

// Int helper function returns integer pointer
func Int(d int) *int {
	return &d
}

// String helper function returns string pointer
func String(s string) *string {
	return &s
}
