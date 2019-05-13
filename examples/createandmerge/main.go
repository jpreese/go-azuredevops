// Copyright 2019 go-azuredevops AUTHORS. All rights reserved.
//
// This example uses the go-azuredevops library to create and then attempt
// to merge a pull request.  You should already have a source and target
// branch available in your repository.  It:
// 0. Prompts the user for required inputs
// 1. Creates a NewClient() using basic auth and personal access token
// 2. Gets the logged in user's IdentityRef{}
// 3. Creates a pull request
// 4. Merges the pull request with the default no fast-forward

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/mcdafydd/go-azuredevops/azuredevops"
	"golang.org/x/crypto/ssh/terminal"
)

type myClient struct {
	c            *azuredevops.Client
	org          string
	project      string
	repo         string
	token        string
	sourceBranch string
	targetBranch string
}

func main() {
	client := new(myClient)

	r := bufio.NewReader(os.Stdin)
	fmt.Print("Azure Devops Organization: ")
	client.org, _ = r.ReadString('\n')
	client.org = strings.TrimSuffix(client.org, "\n")

	r = bufio.NewReader(os.Stdin)
	fmt.Print("Azure Devops Project: ")
	client.project, _ = r.ReadString('\n')
	client.project = strings.TrimSuffix(client.project, "\n")

	r = bufio.NewReader(os.Stdin)
	fmt.Print("Azure Devops Repository: ")
	client.repo, _ = r.ReadString('\n')
	client.repo = strings.TrimSuffix(client.repo, "\n")

	r = bufio.NewReader(os.Stdin)
	fmt.Print("Azure Devops Source Branch Name: ")
	client.sourceBranch, _ = r.ReadString('\n')
	client.sourceBranch = strings.TrimSuffix(client.sourceBranch, "\n")

	r = bufio.NewReader(os.Stdin)
	fmt.Print("Azure Devops Target Branch Name: ")
	client.targetBranch, _ = r.ReadString('\n')
	client.targetBranch = strings.TrimSuffix(client.targetBranch, "\n")

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
	client.create()
}

func (client *myClient) create() *azuredevops.GitPullRequest {
	pull := &azuredevops.GitPullRequest{
		Title:         String("EXAMPLE PULL REQUEST TITLE"),
		Description:   String("EXAMPLE PULL REQUEST DESCRIPTION"),
		SourceRefName: String(client.sourceBranch),
		TargetRefName: String(client.targetBranch),
	}

	pr, _, err := client.c.PullRequests.Create(context.Background(),
		client.org,
		client.project,
		client.repo,
		pull,
	)
	if err != nil {
		fmt.Printf("Error trying to create pull request: %+v\n", err)
		return nil
	}

	fmt.Printf("Successfully created pull request ID %d\nAttempting to merge.\n", pr.GetPullRequestID())
	client.merge(pr, pr.GetCreatedBy(), "Merge using go-azuredevops example")
	return pr
}

func (client *myClient) merge(pull *azuredevops.GitPullRequest, id *azuredevops.IdentityRef, commitMsg string) {
	// Set default pull request completion options
	mcm := azuredevops.NoFastForward.String()
	twi := new(bool)
	*twi = true
	completionOptions := azuredevops.GitPullRequestCompletionOptions{
		BypassPolicy:            new(bool),
		BypassReason:            String(""),
		DeleteSourceBranch:      new(bool),
		MergeCommitMessage:      &commitMsg,
		MergeStrategy:           &mcm,
		SquashMerge:             new(bool),
		TransitionWorkItems:     twi,
		TriggeredByAutoComplete: new(bool),
	}

	merge, _, err := client.c.PullRequests.Merge(context.Background(),
		client.org,
		client.project,
		client.repo,
		pull.GetPullRequestID(),
		pull,
		completionOptions,
		*id,
	)

	if err != nil {
		fmt.Printf("Merge failed: %+v\n", err)
		return
	}
	fmt.Printf("Successfully merged pull request: %+v\n", merge)
}

// String helper function returns string pointer
func String(s string) *string {
	return &s
}
