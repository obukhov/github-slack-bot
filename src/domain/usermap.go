package domain

type UserMap struct {
	userMap map[string]string
}

func NewUserMap(userMap map[string]string) *UserMap {
	return &UserMap{userMap: userMap}
}

func (um *UserMap) GetUserName(user string) string {
	name, found := um.userMap[user]

	if true == found {
		return name
	}

	return user
}

func (um *UserMap) IsDefined(user string) bool {
	_, found := um.userMap[user]

	return found
}
