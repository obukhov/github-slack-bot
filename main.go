package main

import (
	"github.com/obukhov/github-slack-bot/src/domain"
	"github.com/obukhov/github-slack-bot/src/usecase"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"context"
	"log"
	"os"
	"time"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_TOKEN")},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	waitingPRs := make([]domain.WaitingPR, 0)

	workingDir, _ := os.Getwd()
	userMap := usecase.LoadUserMap(workingDir . "/config.yaml")

	// list all repositories for the authenticated user
	pulls, _, err := client.PullRequests.List(
		ctx,
		os.Getenv("GITHUB_OWNER"),
		os.Getenv("GITHUB_REPO"),
		&github.PullRequestListOptions{

			ListOptions: github.ListOptions{PerPage: 100}},
	)

	if nil != err {
		log.Println(err.Error())
	}

	now := time.Now()
	for _, pullRequest := range pulls {

		if false == userMap.IsDefined(pullRequest.User.GetLogin()) {
			continue
		}

		waitingPr := domain.NewWaitingPr(
			pullRequest,
			userMap.GetUserName(pullRequest.User.GetLogin()),
			now.Sub(pullRequest.GetCreatedAt()),
			now.Sub(pullRequest.GetUpdatedAt()),
		)

		for _, user := range pullRequest.Assignees {
			waitingPr.AddReviewStatus(userMap.GetUserName(user.GetLogin()), "assigned")
		}

		waitingPRs = append(waitingPRs, *waitingPr)
	}

	log.Println(waitingPRs)
}
