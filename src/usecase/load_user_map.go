package usecase

import (
	"github.com/obukhov/github-slack-bot/src/domain"
	"github.com/go-yaml/yaml"
)

func LoadUserMap(path string) domain.UserMap {
	rawMap := make(map[string]interface{})

	yaml.Unmarshal()
}
