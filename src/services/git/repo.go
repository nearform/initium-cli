package git

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

func PublishCommentPRGithub (url string) error {
	var message, owner, repo string
	var prNumber int
	commitSha, err := GetHash()

	// Build message
	message = fmt.Sprintf("Application URL: %s\n", url) + fmt.Sprintf("Commit hash: %s\n", commitSha) + fmt.Sprintf("Timestamp: %v\n", time.Now())

	// Debug
	fmt.Println(message)

	// Check GITHUB_TOKEN
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("Please set up the GITHUB_TOKEN environment variable")
	}

	// Create an authenticated GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get required data to publish a comment
	repoInfo := os.Getenv("GITHUB_REPOSITORY")
	repoParts := strings.Split(repoInfo, "/")
	if len(repoParts) == 2 {
		owner = repoParts[0]
		repo = repoParts[1]

		fmt.Printf("Owner: %s\n", owner) // Debug
		fmt.Printf("Repository: %s\n", repo) // Debug
	} else {
		return fmt.Errorf("Invalid repository information")
	}

	// Check if the workflow was triggered by a pull request event
	eventName := os.Getenv("GITHUB_EVENT_NAME")
	if eventName == "pull_request" {
		// Get the pull request ref
		prRef := os.Getenv("GITHUB_REF")

		// Extract the pull request number using a regular expression
		re := regexp.MustCompile(`refs/pull/(\d+)/merge`)
		matches := re.FindStringSubmatch(prRef)

		if len(matches) == 2 {
			prNumber, err = strconv.Atoi(matches[1])
			if err != nil {
				return fmt.Errorf("Error converting string to int: %v", err)
			}
			fmt.Printf("Pull Request Number: %d\n", prNumber) // Debug
		} else {
			return fmt.Errorf("Unable to extract pull request number from GITHUB_REF")
		}
	} else {
		return fmt.Errorf("This workflow was not triggered by a pull request event")
	}

	// Create comment with body
	comment := &github.IssueComment{
		Body: github.String(message),
	}

	// List comments on the PR
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, prNumber, nil)
	if err != nil {
		fmt.Printf("Error listing comments: %v\n", err) // Debug
		return err
	}
	commentID := findExistingCommentIDPRGithub(comments, fmt.Sprintf("Application URL: %s\n", url)) // Search for app URL comment

	if commentID != 0 {
		// Update existing comment
		updatedComment, _, err := client.Issues.EditComment(ctx, owner, repo, commentID, comment)
		if err != nil {
			fmt.Printf("Error updating comment: %v\n", err) // Debug
			return err
		}
		fmt.Printf("Comment updated successfully: %s\n", updatedComment.GetHTMLURL())
	} else {
		// Publish a new comment
		newComment, _, err := client.Issues.CreateComment(ctx, owner, repo, prNumber, comment)
		if err != nil {
			fmt.Printf("Error publishing new comment: %v\n", err) // Debug
			return err
		}
		fmt.Printf("Comment published: %s\n", newComment.GetHTMLURL())
	}

	return nil
}

func findExistingCommentIDPRGithub(comments []*github.IssueComment, targetBody string) int64 {
	for _, comment := range comments {
		if strings.TrimSpace(comment.GetBody()) == strings.TrimSpace(targetBody) {
			return comment.GetID()
		}
	}
	return 0
}
