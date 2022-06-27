package minioHandler

type Config struct {
	Workers         int
	ChannelCapacity int
	BucketName      string
	Endpoint        string
	AccessKeyID     string
	SecretKeyID     string
	UseSSL          bool
}

func NewConfig() *Config {
	return &Config{
		BucketName:      "files",
		Workers:         1,
		ChannelCapacity: 1,
		UseSSL:          false,
		Endpoint:        "localhost:9000",
	}
}
