package config

type Config struct {
	// Credentials
	Credentials struct {
		Username string
		Domain   string
		Password string
		Hashes   string
		DCIP     string
	}
	// Network
	UseLdaps      bool
	DnsNameServer string
	// General
	Threads   int
	OutputDir string
	Debug     bool
}
