package server

import "strings"

type Config struct {
	LocalDev bool `env:"LOCAL_DEV" envDefault:"false"`

	Port int `env:"BLOG_PORT" envDefault:"9191"`

	DBUser     string `env:"DATABASE_USER"`
	DBPassword string `env:"DATABASE_PASSWORD"`
	DBHostname string `env:"DATABASE_HOST"`
	DBPort     int    `env:"DATABASE_PORT" envDefault:"3306"`
	DBName     string `env:"DATABASE_DB"   envDefault:"blog3"`

	AdminUser     string `env:"ADMIN_USER"   envDefault:"admin"`
	AdminPassword string `env:"ADMIN_PW"   envDefault:"admin"`

	HubUrls string `env:"HUB_URLS"`

	// 9*60*60=32400 is JST
	TimeZoneOffset int `env:"TIMEZONE_OFFSET" envDefault:"32400"`
}

func (c *Config) GetHubUrls() []string {
	if c.HubUrls != "" {
		return strings.Split(c.HubUrls, ",")
	}
	return []string{
		"https://pubsubhubbub.appspot.com/",
		"https://pubsubhubbub.superfeedr.com/",
	}
}
