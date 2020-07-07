package utils

type User struct {
	Username     string
	Key          string
	BaseFilePath string
	ConfigPath   string
	Objects      map[string]FObject
}

type UserRequestPackage struct {
	ClientUser User
	Object     FObject
}
