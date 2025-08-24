package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	ihttp "github.com/dya-andrade/cat-api/internal/http"
	"github.com/dya-andrade/cat-api/internal/service"
	"github.com/dya-andrade/cat-api/internal/storage"
	"github.com/dya-andrade/cat-api/internal/worker"

	"github.com/dya-andrade/cat-api/internal/config"
)

func main() {
	// Cria um contexto que escuta sinais do sistema (Ctrl+C, kill, etc) para fazer shutdown gracioso do servidor
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop() // Garante que o contexto será finalizado ao sair da main

	// Carrega configurações do sistema (porta, banco, timeout, etc)
	cfg := config.MustLoad()
	log.Printf("starting cats-api on %s", cfg.AppAddr)

	// Conecta ao banco de dados Postgres usando as configs
	pg, err := storage.NewPostgres(ctx, cfg.DBDsn, cfg.DBMaxConns, cfg.DBMinConns, cfg.DBMaxIdleTime)
	if err != nil {
		log.Fatalf("db connect error: %v", err) // Encerra o programa se não conectar
	}
	defer pg.Close() // Fecha a conexão ao sair

	// Cria o repositório de gatos usando o pool de conexões do banco
	catRepo := storage.NewCatRepository(pg.Pool)

	// Cria um pool de workers para tarefas assíncronas (ex: gerar thumbnails)
	wp := worker.NewPool(int(cfg.WorkerConcurrency))
	wp.Start()          // Inicia os workers
	defer wp.Shutdown() // Garante que os workers serão finalizados ao sair

	// Cria o serviço de gatos, passando repositório, pool de workers e timeout
	catSvc := service.NewCatService(catRepo, wp, cfg.RequestTimeout)

	// Cria o roteador HTTP e configura o servidor
	router := ihttp.NewRouter(catSvc)
	srv := &http.Server{
		Addr:         cfg.AppAddr,      // Endereço e porta do servidor
		Handler:      router,           // Handler das rotas
		ReadTimeout:  15 * time.Second, // Timeout para leitura da requisição
		WriteTimeout: 15 * time.Second, // Timeout para escrita da resposta
		IdleTimeout:  60 * time.Second, // Timeout para conexões ociosas
	}

	// Inicia o servidor HTTP em uma goroutine e aguarda erro ou sinal de shutdown
	errCh := make(chan error, 1)
	go func() {
		log.Printf("http server listening at http://%s", cfg.AppAddr)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err // Envia erro para o canal se não for shutdown normal
		}
	}()

	// Espera por sinal de encerramento ou erro do servidor
	select {
	case <-ctx.Done():
		log.Println("shutdown signal received") // Recebeu sinal do sistema
	case err := <-errCh:
		log.Printf("server error: %v", err) // Ocorreu erro no servidor
	}

	// Faz shutdown gracioso do servidor HTTP, aguardando até 10 segundos para finalizar requisições pendentes
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}

	// Finaliza o pool de workers (drena fila de tarefas)
	wp.Shutdown()
	log.Println("bye!")
	_ = os.Stderr // Evita erro de import não usado em alguns ambientes
}
