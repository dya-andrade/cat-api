package storage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// nil -> nulo

// Postgres is a PostgreSQL database connection pool.
type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, dsn string, maxConns, minConns int32, maxIdleTime time.Duration) (*Postgres, error) {
	// Recebe contexto, string de conexão, limites de conexões e tempo máximo ocioso

	cfg, err := pgxpool.ParseConfig(dsn)
	// Cria uma configuração de pool de conexões a partir da string de conexão (DSN)
	if err != nil {
		return nil, err
		// Se a string de conexão for inválida, retorna erro
	}

	cfg.MaxConns = maxConns
	// Define o número máximo de conexões simultâneas no pool

	cfg.MinConns = minConns
	// Define o número mínimo de conexões abertas no pool

	cfg.MaxConnIdleTime = maxIdleTime
	// Define o tempo máximo que uma conexão pode ficar ociosa antes de ser fechada

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	// Cria o pool de conexões usando a configuração e o contexto
	if err != nil {
		return nil, err
		// Se houver erro ao criar o pool, retorna erro
	}

	return &Postgres{Pool: pool}, nil
	// Retorna uma instância de Postgres com o pool criado e nil para erro
}

// Close fecha o pool de conexões com o banco de dados.
func (p *Postgres) Close() {
	p.Pool.Close()
}
