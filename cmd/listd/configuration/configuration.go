package configuration

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var (
	// ErrEnvNoDefault is an error that denotes an environment variable was not supplied
	// that has no default value
	ErrEnvNoDefault = errors.New("env variable not supplied and has no default")
)

// Config is the struct that contains fields that store the necessary configuration
// gathered from the environment
type Config struct {
	DaemonPort int

	DBUser string
	DBPass string
	DBName string
	DBHost string
	DBPort int

	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

// Environment attempts to gather all of the environment variables from the context of where
// the program is ran.
func Environment() (*Config, error) {
	var (
		c   Config
		err error
	)

	if c.DaemonPort, err = envInt(EnvDaemonPort, DefaultDaemonPort); err != nil {
		return nil, errors.Wrap(err, "get daemon port")
	}

	if c.DBUser = os.Getenv(EnvDBUser); c.DBUser == "" {
		c.DBUser = DefaultDBUser
	}

	if c.DBPass = os.Getenv(EnvDBPass); c.DBPass == "" {
		c.DBPass = DefaultDBPass
	}

	if c.DBName = os.Getenv(EnvDBName); c.DBName == "" {
		c.DBName = DefaultDBName
	}

	if c.DBHost = os.Getenv(EnvDBHost); c.DBHost == "" {
		c.DBHost = DefaultDBHost
	}

	if c.DBPort, err = envInt(EnvDBPort, DefaultDBPort); err != nil {
		return nil, errors.Wrap(err, "get postgres port")
	}

	var second int
	if second, err = envInt(EnvReadTimeout, DefaultReadTimeout); err != nil {
		return nil, errors.Wrap(err, "get read timeout")
	}
	c.ReadTimeout = time.Second * time.Duration(second)

	if second, err = envInt(EnvWriteTimeout, DefaultWriteTimeout); err != nil {
		return nil, errors.Wrap(err, "get write timeout")
	}
	c.WriteTimeout = time.Second * time.Duration(second)

	if second, err = envInt(EnvShutdownTimeout, DefaultShutdownTimeout); err != nil {
		return nil, errors.Wrap(err, "get shutdown timeout")
	}
	c.ShutdownTimeout = time.Second * time.Duration(second)

	return &c, nil
}

// envInt is a utility function that returns the integer stored at
// the given environment variable. -1 default value denotes that there
// is no default value.
func envInt(name string, defaultValue int) (int, error) {
	var s string
	if s = os.Getenv(name); s == "" {
		if defaultValue == -1 {
			return 0, ErrEnvNoDefault
		}

		return defaultValue, nil
	}

	return strconv.Atoi(s)
}
