package configuration

const (
	// EnvDaemonPort is the name of the environment variable that stores the
	// port that the daemon runs on
	EnvDaemonPort     = "LIST_DAEMON_PORT"
	DefaultDaemonPort = 3000

	// EnvDBUser is the name of the environment variable that stores the postgres
	// database username
	EnvDBUser     = "DB_USER"
	DefaultDBUser = "root"

	// EnvDBPass is the name of the environment variable that stores the postgres
	// database password
	EnvDBPass     = "DB_PASS"
	DefaultDBPass = "root"

	// EnvDBName is the name of the environment variable that stores the postgres
	// database name
	EnvDBName     = "DB_NAME"
	DefaultDBName = "list"

	// EnvDBHost is the name of the environment variable that stores the postgres
	// database host name
	EnvDBHost     = "DB_HOST"
	DefaultDBHost = "db"

	// EnvDBPort is the name of the environment variable that stores the postgres
	// database port number
	EnvDBPort     = "DB_PORT"
	DefaultDBPort = 5432

	// EnvReadTimeout is the name of the environment variable that stores the
	// time in seconds for the timeout of reading actions for the server
	EnvReadTimeout     = "READ_TIMEOUT"
	DefaultReadTimeout = 5

	// EnvWriteTimeout is the name of the environment variable that stores the
	// time in seconds for the timeout of writing actions for the server
	EnvWriteTimeout     = "WRITE_TIMEOUT"
	DefaultWriteTimeout = 10

	// EnvWriteTimeout is the name of the environment variable that stores the
	// time in seconds for the timeout in regards to gracefully shutting down
	// the server
	EnvShutdownTimeout     = "SHUTDOWN_TIMEOUT"
	DefaultShutdownTimeout = 5
)
