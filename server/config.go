package server

const (
	defaultHttpPort = 8080
)

// Config is the server configuration. This enables us to change execution environment of the server
// without creating a new build; keeps us away from hard-coding stuff.
type Config struct {
	HttpPort    int
	DatabaseUrl string
}

// NewConfig returns a Config with sensible default values assigned to some fields.
func NewConfig() *Config {
	return &Config{
		HttpPort: defaultHttpPort,
	}
}
