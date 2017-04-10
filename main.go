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
	"github.com/dustin/go-humanize"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_API_TOKEN")},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	waitingPRsByList := make(map[string][]domain.WaitingPR)

	workingDir, _ := os.Getwd()
	userMap, _ := usecase.LoadUserMap(workingDir + "/config.yml")

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

		if false == userMap.HasGithubUser(pullRequest.User.GetLogin()) {
			continue
		}

		slackUserName := userMap.SlackUserName(pullRequest.User.GetLogin())
		waitingPr := domain.NewWaitingPr(
			pullRequest,
			slackUserName,
			now.Sub(pullRequest.GetCreatedAt()),
			now.Sub(pullRequest.GetUpdatedAt()),
		)

		for _, user := range pullRequest.Assignees {
			waitingPr.AddReviewStatus(userMap.SlackUserName(user.GetLogin()), "assigned")
		}

		channel := userMap.Channel(slackUserName)

		if _, found := waitingPRsByList[channel]; false == found {
			waitingPRsByList[channel] = make([]domain.WaitingPR, 0)
		}

		waitingPRsByList[channel] = append(waitingPRsByList[channel], *waitingPr)

	}

	for teamName, waitingPrList := range waitingPRsByList {
		log.Printf("For team %s:", teamName)
		for _, pr := range waitingPrList {
			log.Printf(
				"Pull request [%d] %s by @%s is waiting for %s",
				pr.Pr.Number,
				pr.Pr.GetTitle(),
				pr.Author,
				humanize.RelTime(time.Now(), time.Now().Add(pr.WaitingSinceCreated), "", ""),
			)
		}

	}
}
