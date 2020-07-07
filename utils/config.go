package utils

type ConfigCloudStore struct {
	ServerDbPath   string
	ServerBasePath string `default:"/Users/harshilpatel/Projects/test/cloud-store-server"`
}

func GetConfiguration() ConfigCloudStore {
	return ConfigCloudStore{
		ServerDbPath:   "/Users/harshilpatel/Projects/test/cloud-store-server",
		ServerBasePath: "/Users/harshilpatel/Projects/test/cloud-store-server",
	}
}
