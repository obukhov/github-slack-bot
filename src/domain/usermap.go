package domain

import "errors"

type UserMap struct {
	githubUserToSlackUser map[string]string
	slackUserToChannel    map[string]string
}

func NewUserMap() *UserMap {
	return &UserMap{
		githubUserToSlackUser: make(map[string]string),
		slackUserToChannel:    make(map[string]string),
	}
}

func (um *UserMap) AddUserTeam(channelName, githubUserName, slackUserName string) error {
	if um.HasGithubUser(githubUserName) {
		return errors.New("User already exist")
	}

	um.githubUserToSlackUser[githubUserName] = slackUserName
	um.slackUserToChannel[slackUserName] = channelName

	return nil
}

func (um *UserMap) HasGithubUser(user string) bool {
	_, ok := um.githubUserToSlackUser[user]
	return ok
}

func (um *UserMap) SlackUserName(githubUserName string) string {
	return um.githubUserToSlackUser[githubUserName]
}

func (um *UserMap) Channel(slackUserName string) string {
	return um.slackUserToChannel[slackUserName]
}
