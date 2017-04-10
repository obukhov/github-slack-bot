package usecase

import (
	"github.com/obukhov/github-slack-bot/src/domain"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
)

type User struct {
	GitHubName string `yaml:"github"`
	SlackName  string `yaml:"slack"`
}

type Team struct {
	Name  string `yaml:"team"`
	Users []User `yaml:"users"`
}

type Config struct {
	UserMap []Team `yaml:"user_map"`
}

func LoadUserMap(path string) (*domain.UserMap, error) {
	log.Println(path)

	data, err := ioutil.ReadFile(path)
	if nil != err {
		log.Println(err)
		return nil, err
	}

	config := &Config{}
	yaml.Unmarshal(data, &config)

	userMap := domain.NewUserMap()

	for _, team := range config.UserMap {
		for _, user := range team.Users {
			userMap.AddUserTeam(team.Name, user.GitHubName, user.SlackName)
		}
	}

	return userMap, nil
}
