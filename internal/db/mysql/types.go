package mysql

type Config struct {
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	Charset   string
}

const DefaultCharset = "utf8mb4"
const DefaultPort = 3306