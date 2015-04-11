package lib

import (
   "github.com/ccding/go-logging/logging"
   "github.com/danryan/env"
   "os"
   "time"
)

// Config struct
type Config struct {
   Name        string `env:"key=LAZLO_NAME default=lazlo"`
   Token 	   string `env:"key=LAZLO_TOKEN"`
   URL	 	   string `env:"key=LAZLO_URL default=http://localhost"`
   LogLevel    string `env:"key=LAZLO_LOG_LEVEL default=info"`
   RedisURL 	string `env:"key=LAZLO_REDIS_URL"`
   RedisPW 		string `env:"key=LAZLO_REDIS_PW"`
   Port 			string `env:"key=PORT"`
}

func newConfig() *Config {
   c := &Config{}
   env.MustProcess(c)
   return c
}

func newLogger() *logging.Logger {
   format := "%25s [%s] %8s: %s\n time,name,levelname,message"
   timeFormat := time.RFC3339
   level := logging.GetLevelValue(`INFO`)
   logger, _ := logging.WriterLogger("lazlo", level, format, timeFormat, os.Stdout, true)
   return logger
}
