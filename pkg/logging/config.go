package logging

type Config struct {
	IsJson bool `env:"SHINFO_LOGGING_JSON, default=false"`
	Level  int8 `env:"SHINFO_LOGGING_LEVEL, default=0"`
}
