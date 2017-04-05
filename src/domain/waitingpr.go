package domain

import (
	"github.com/google/go-github/github"
	"time"
)

const (
	STATUS_APPROVED        = "approved"
	STATUS_COMMENTED       = "commented"
	STATUS_REQUEST_CHANGES = "request_changes"
)

type WaitingPR struct {
	Pr                     *github.PullRequest
	Author                 string
	ReviewStatus           map[string]string
	WaitingSinceCreated    time.Duration
	WaitingSinceLastChange time.Duration
}

func NewWaitingPr(request *github.PullRequest, author string, sinceCreated, sinceUpdated time.Duration) *WaitingPR {
	return &WaitingPR{
		request,
		author,
		make(map[string]string),
		sinceCreated,
		sinceUpdated,
	}
}

func (wpr *WaitingPR) AddReviewStatus(username string, status string) {
	wpr.ReviewStatus[username] = status
}
