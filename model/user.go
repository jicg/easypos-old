package model

func GetUserByNameWithPwd(name string, pwd string) *User {
	user := new(User)
	Engine.Where(" name = ? ", name).And(" pwd = ? ", pwd).Get(user)
	return user
}

func GetUserCount() (int64, error) {
	return Engine.Count(new(User))
}
