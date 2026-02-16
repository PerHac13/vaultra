package s3

type Config struct {
	Bucket            string
	Region            string
	AccessKeyID       string
	SecretAccessKey   string
	Prefix            string
	Endpoint          string
	UsePathStyle      bool
	DisableSSL        bool
	PartSize          int64
}

const DefaultPartSize = 5 * 1024 * 1024 // 5 MB
const DefaultRegion = "us-east-1"
const DefaultPrefix = "vaultra/backups/"

