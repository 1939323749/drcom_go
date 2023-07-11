package conf

import (
	"flag"
	"os"
)

// Conf global variable.
var (
	Conf     *Config
	confPath string

	version    = flag.String("version", "1.0", "version")
	authServer = flag.String("auth_server", "10.100.61.3", "auth server")
	port       = flag.String("port", "61440", "port")
	username   = flag.String("username", "USERNAME", "username")
	password   = flag.String("password", "PASSWORD", "password")
	hostname   string
	mac        = flag.String("mac", "00:00:00:00:00:00", "MAC address")
	ip         = flag.String("ip", "0.0.0.0", "IP address")
)

type Config struct {
	Version    string
	AuthServer string
	Port       string
	Username   string
	Password   string
	Hostname   string
	MAC        string
	IP         string
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init create config instance.
func Init() (err error) {
	Conf = &Config{
		Version:    *version,
		AuthServer: *authServer,
		Port:       *port,
		Username:   *username,
		Password:   *password,
		MAC:        *mac,
		IP:         *ip,
	}

	if hostname, err = os.Hostname(); err != nil {
		hostname = "unknown"
	}
	Conf.Hostname = hostname
	return
}
