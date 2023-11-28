package git

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	git "github.com/go-git/go-git/v5"
	github "github.com/google/go-github/v56/github"
	oauth2 "golang.org/x/oauth2"
)

const (
	httpgithubprefix        = "https://github.com/"
	gitgithubprefix         = "git@github.com:"
	githubRefPattern        = `refs/pull/(\d+)/merge`
	githubCommentIdentifier = "Deployed by [Initium](https://initium.nearform.com)" // used to find and update existing comments
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

func buildMarkdownMessage(url string) (string, error) {
	commitSha, err := GetHash()
	if err != nil {
		return "", err
	}

	message := fmt.Sprintf(githubCommentIdentifier+`
|Application URL | %s |
|:-----------------|:----|
|Commit hash | %s |
|Timestamp | %s |
`, url, commitSha, time.Now().UTC())
	return message, nil
}

func PublishCommentPRGithub(url string) error {
	token := os.Getenv("GITHUB_TOKEN")
	prRef := os.Getenv("GITHUB_REF")
	repoInfo := os.Getenv("GITHUB_REPOSITORY")

	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable not set")
	}

	// Extract pull request number
	prNumber, err := extractPullRequestNumber(prRef)
	if err != nil {
		return err
	}

	message, err := buildMarkdownMessage(url)
	if err != nil {
		return fmt.Errorf("cannot build the message: %v", err)
	}

	comment := &github.IssueComment{
		Body: github.String(message),
	}

	// Get required data to publish a comment
	repoParts := strings.Split(repoInfo, "/")
	if len(repoParts) != 2 {
		return fmt.Errorf("invalid repository information %s", repoInfo)
	}
	owner := repoParts[0]
	repo := repoParts[1]

	// Create an authenticated GitHub client
	ctx := context.Background()
	client := createGithubClient(ctx, token)

	// Check if we have to update an existing comment
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, prNumber, nil)
	if err != nil {
		return err
	}

	matchingComments := findExistingGithubComments(comments, githubCommentIdentifier) // Search for app URL comment
	if n := len(matchingComments); n != 0 {
		log.Infof("%d matching comment[s] found %v, will always update the last one", n, matchingComments)
		updatedComment, _, err := client.Issues.EditComment(ctx, owner, repo, matchingComments[n-1], comment)
		if err != nil {
			return err
		}
		log.Infof("Comment updated successfully: %s\n", updatedComment.GetHTMLURL())
		return nil
	}

	// Publish a new comment
	newComment, _, err := client.Issues.CreateComment(ctx, owner, repo, prNumber, comment)
	if err != nil {
		return err
	}
	log.Infof("Comment published: %s\n", newComment.GetHTMLURL())
	return nil
}

func createGithubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func extractPullRequestNumber(prRef string) (int, error) {
	matches := regexp.MustCompile(githubRefPattern).FindStringSubmatch(prRef)
	if len(matches) != 2 {
		return 0, fmt.Errorf("unable to extract pull request number from GITHUB_REF %s", prRef)
	}

	prNumber, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("error converting string to int: %v", err)
	}
	return prNumber, nil
}

func findExistingGithubComments(comments []*github.IssueComment, targetString string) []int64 {
	matchingComments := []int64{}
	for _, comment := range comments {
		body := comment.GetBody()
		if strings.Contains(body, targetString) && strings.Contains(body, "initium") {
			matchingComments = append(matchingComments, comment.GetID())
		}
	}
	return matchingComments
}
