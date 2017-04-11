package main

import (
	"github.com/obukhov/github-slack-bot/src/domain"
	"github.com/obukhov/github-slack-bot/src/usecase"
	"context"
	"log"
	"os"
	"time"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"github.com/dustin/go-humanize"
	"golang.org/x/oauth2"
	"github.com/lucasb-eyer/go-colorful"
	"strings"
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
			ListOptions: github.ListOptions{PerPage: 100},
		},
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
			pullRequest.User.GetLogin(),
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

	apiClient := slack.New(os.Getenv("SLACK_TOKEN"))
	for teamName, waitingPrList := range waitingPRsByList {
		attachments := make([]slack.Attachment, 0)
		for _, pr := range waitingPrList {
			assigneeNames := make([]string, 0)
			for _, assignee := range pr.Pr.Assignees {
				assigneeNames = append(assigneeNames, "@"+userMap.SlackUserName(assignee.GetLogin()))
			}

			assignees := ""
			if len(assigneeNames) > 0 {
				assignees = fmt.Sprintf("assigned to %s ", strings.Join(assigneeNames, ", "))
			}

			waitingFor := humanize.RelTime(now, now.Add(pr.WaitingSinceCreated), "", "")

			attachments = append(attachments, slack.Attachment{
				Title:     fmt.Sprintf(pr.Pr.GetTitle()),
				TitleLink: pr.Pr.GetHTMLURL(),
				Color:     colorful.HappyColor().Hex(),
				ThumbURL:  pr.Pr.User.GetAvatarURL(),
				Text: fmt.Sprintf(
					"by %s %sis waiting for %s",
					pr.Author,
					assignees,
					waitingFor,
				),
			})

		}

		pullReqCount := len(attachments)
		if pullReqCount > 0 {
			_, _, err := apiClient.PostMessage(
				teamName,
				fmt.Sprintf("You have %d pull requests unattended", pullReqCount),
				slack.PostMessageParameters{
					Username:    "github",
					IconEmoji:   ":cat:",
					Parse:       "full",
					Attachments: attachments,
				});
			if nil == err {
				log.Printf("Post message to channel %s about %d pull requests", teamName, pullReqCount)
			} else {
				log.Printf("Error posting message about %d pull requests to channel %s: %s", pullReqCount, teamName, err.Error())
			}
		} else {
			log.Printf("No pull requests for team %s", teamName)

		}

	}

}
