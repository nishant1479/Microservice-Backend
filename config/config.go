package config

import "time"

const (
	GRPC_PORT = "GRPC_PORT"
	HTTP_PORT = "HTTP_PORT"
)

// Configuration of the application
type Config struct {
	AppVersion string
	Server     Server
	Logger     Logger
	Jeager     Jeager
	Metrics    Metrics
	MongoDB    MongoDB
	Kafka      Kafka
	Http       Http
	Redis      Redis
}

// Configuration of the server
type Server struct {
	Port        		string
	Development 		bool
	Timeout     		time.Duration
	ReadTimeout			time.Duration
	WriteTimeout		time.Duration
	MaxConnectionIdle	time.Duration
	MaxConnectionAge	time.Duration
	Kafka				Kafka
}


// Configuration of the Http
type Http struct{
	Port				string
	PprofPort			string
	Timeout				time.Duration
	ReadTimeout			time.Duration
	WriteTimeout		time.Duration
	CookieLifeTime		int
	SessionCookieName	string
}


// Configuration of the Logger
type Logger struct{
	DisableCaller		bool
	DisableStackTrace	bool
	Encoding			string
	Level				string
}


// Configuration of the Metrics
type Metrics struct{
	Port		string
	URL			string
	ServiceName	string
}

type Jeager struct{
	Host		string
	ServiceName	string
	LogSpans	bool
}

type MongoDB struct{
	URI			string
	User		string
	Password	string
	DB			string
}

type Kafka struct{
	Brokers []string
}

type Redis struct{
	RedisAddr		string
	RedisPassword	string
	RedisDB			string
	RedisDefaultDB	string
	MinIdleConn		int
	PoolSize		int
	PoolTimeout		int
	Password		string
	DB				int
}


