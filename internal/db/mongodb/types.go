package mongodb

import "fmt"

type Config struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	AuthDatabase    string
	URI			    string
}

const DefaultPort = 27017
const DefaultAuthDatabase = "admin"

func (c *Config) GetURI() string {
	if c.URI != "" {
		return c.URI
	}
	if c.Username != "" && c.Password != "" {
		return "mongodb://" + c.Username + ":" + c.Password + "@" + c.Host + ":" + fmt.Sprintf("%d",c.Port) + "/" + c.Database
	}
	return "mongodb://" + c.Host + ":" + fmt.Sprintf("%d",c.Port)
}