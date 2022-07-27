package models

// Config defines the configuration for a single server
type Config struct {
	// Name of the server
	Name string
	// hostname/ip:port
	Host string
	// The type of the server for when we have master and leaf servers
	Type string
	// URL for the outside access to api
	URL string
	// cli flag to turn off redis support
	NoRedis bool
}
