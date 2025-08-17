package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppAddr           string        // Endereço e porta onde a aplicação vai rodar (ex: ":8080")
	DBDsn             string        // String de conexão do banco de dados (Data Source Name)
	DBMaxConns        int32         // Número máximo de conexões simultâneas no banco
	DBMinConns        int32         // Número mínimo de conexões abertas no banco
	DBMaxIdleTime     time.Duration // Tempo máximo que uma conexão pode ficar ociosa
	WorkerConcurrency int32         // Quantidade de workers para processar tarefas em paralelo
	RequestTimeout    time.Duration // Tempo limite para cada requisição HTTP
}

func getEnvString(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}

func getEnvInt32(key, def string) int32 {
	val, _ := strconv.Atoi(getEnvString(key, def))
	return int32(val)
}

func getEnvDuration(key, def string) time.Duration {
	val, _ := time.ParseDuration(getEnvString(key, def))
	return val
}

func MustLoad() Config {
	return Config{
		AppAddr:           getEnvString("APP_ADDR", ":8080"),
		DBDsn:             getEnvString("DB_DSN", "postgres://cat_user:cat_password@localhost:5432/cat_db?sslmode=disable"),
		DBMaxConns:        getEnvInt32("DB_MAX_CONNS", "10"),
		DBMinConns:        getEnvInt32("DB_MIN_CONNS", "2"),
		DBMaxIdleTime:     getEnvDuration("DB_MAX_IDLE_TIME", "30s"),
		WorkerConcurrency: getEnvInt32("WORKER_CONCURRENCY", "4"),
		RequestTimeout:    getEnvDuration("REQUEST_TIMEOUT", "10s"),
	}
}
