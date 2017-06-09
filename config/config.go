package config

type Config struct {
	Database Database
	Trace    Trace
	Debug    bool
	Httpport string
}

type Trace struct {
	Zipkin struct {
		Collector struct{ Url string }
		Recoder   struct{ HostPort, ServiceName string }
	}
}

type Database struct {
	Driver, Connection string
	Logger             struct {
		Kafka struct {
			Brokers []string
			Topic   string
		}
	}
}
