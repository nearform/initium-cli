package git

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	github "github.com/google/go-github/v56/github"
	oauth2 "golang.org/x/oauth2"
)

const (
	httpgithubprefix = "https://github.com/"
	gitgithubprefix  = "git@github.com:"
)

func initRepo() (*git.Repository, error) {
	wd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	r, err := git.PlainOpen(wd) // path to the repository
	if err != nil {
		return nil, err
	}

	return r, nil
}

func GetHash() (string, error) {
	repo, err := initRepo()

	if err != nil {
		return "", fmt.Errorf("error getting git repo %v", err)
	}

	headRef, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("cannot get HEAD reference: %v", err)
	}

	return headRef.Hash().String(), nil
}

func GetBranchName() (string, error) {
	repo, err := initRepo()

	if err != nil {
		return "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	branchName := head.Name().Short()
	return branchName, nil
}

func getGithubRemote() (string, error) {
	repo, err := initRepo()

	if err != nil {
		return "", fmt.Errorf("error getting git repo %v", err)
	}

	c, err := repo.Config()
	if err != nil {
		return "", fmt.Errorf("error getting the repo config: %v", err)
	}

	for _, remote := range c.Remotes {
		for _, url := range remote.URLs {
			// HTTPS (https://github.com/organization/repo.git)
			if strings.HasPrefix(url, httpgithubprefix) {
				return strings.Replace(url, httpgithubprefix, "", 1), nil
			}
			// SSH (git@github.com:organization/repo.git)
			if strings.HasPrefix(url, gitgithubprefix) {
				return strings.Replace(url, gitgithubprefix, "", 1), nil
			}
		}
	}

	return "", fmt.Errorf("no github remote found")
}

func GetRepoName() (string, error) {
	remote, err := getGithubRemote()

	if err != nil {
		return "", err
	}

	splitRemote := strings.Split(remote, "/")
	return strings.Replace(splitRemote[1], ".git", "", 1), nil
}

func GetGithubOrg() (string, error) {
	remote, err := getGithubRemote()

	if err != nil {
		return "", err
	}

	splitRemote := strings.Split(remote, "/")
	return splitRemote[0], nil
}

func PublishCommentPRGithub (url string) {
	var message string
	commitSha, err := GetHash()

	// Build message
	message = fmt.Sprintf("Application URL: %s\n", url) + fmt.Sprintf("Commit hash: %s\n", commitSha) + fmt.Sprintf("Timestamp: %v\n", time.Now())

	// Debug
	fmt.Println(message)

	// Check GITHUB_TOKEN
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Printf("Please set your GitHub access token in the GITHUB_TOKEN environment variable.")
	}

	// Create an authenticated GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// TODO: Replace these variables with your repository owner, repository name, and pull request number
	owner := "mablanco"
	repo := "initium-nodejs-demo-app"
	prNumber := 2

	// Specify the comment body
	comment := &github.IssueComment{
		Body: github.String(message),
	}

	// Post the comment to the pull request
	newComment, _, err := client.Issues.CreateComment(ctx, owner, repo, prNumber, comment)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Printf("Comment created: %s\n", newComment.GetHTMLURL())
}
