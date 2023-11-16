package git

import (
	"context"
	"fmt"
	"log"
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
	commitSha, err := GetHash()

	// Debug
	fmt.Printf("Application URL: %s", url)
	fmt.Printf("Commit hash: %s", commitSha)
	fmt.Printf("Timestamp: %v", time.Now())

	// Check GITHUB_TOKEN
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Please set your GitHub access token in the GITHUB_TOKEN environment variable.")
	}

	// Create an authenticated GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Replace these variables with your repository owner, repository name, and pull request number
	owner := "<your username>"
	repo := "<your repository>"
	prNumber := 1

	// Specify the comment body
	comment := &github.PullRequestComment{
		Body: github.String("<your comment here>"),
	}

	// Post the comment to the pull request
	newComment, _, err := client.PullRequests.CreateComment(ctx, owner, repo, prNumber, comment)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Comment created: %s\n", newComment.GetHTMLURL())
}
