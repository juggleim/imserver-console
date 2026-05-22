package models

const (
	ChannelQiNiu = "qiniu"
	ChannelAws   = "aws"
	ChannelOss   = "oss"
	ChannelMinio = "minio"
)

type S3Config struct {
	AccessKey string `json:"access_key,omitempty"`
	SecretKey string `json:"secret_key,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"`
	Region    string `json:"region,omitempty"`
	Bucket    string `json:"bucket,omitempty"`
}

type QiNiuConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	Domain    string `json:"domain"`
}

type OssConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Endpoint  string `json:"endpoint"`
	Bucket    string `json:"bucket"`
	Region    string `json:"region"`
	Domain    string `json:"domain"`
}

type MinioConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Endpoint  string `json:"endpoint"`
	UseSSL    bool   `json:"use_ssl"`
	Bucket    string `json:"bucket"`
}
