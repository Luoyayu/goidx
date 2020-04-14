package config

var (
	WorkerHost = "???.???.workers.dev"
	UserName   = "?????"
	Password   = "?????"
)

var (
	RootBasicAuth = "Basic "
)

type Config struct {
	Host     string `toml:"host"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}
