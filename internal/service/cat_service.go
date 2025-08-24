package service

import (
	"context"
	"errors"
	"time"

	"github.com/dya-andrade/cat-api/internal/domain"
	"github.com/dya-andrade/cat-api/internal/worker"
)

var (
	ErrNotFound = errors.New("cat not found")
)

// CatRepository descreve o que o serviço precisa do repositório.
// Define as operações que o serviço pode executar no banco de dados.
type CatRepository interface {
	Create(ctx context.Context, in domain.CatCreate) (domain.Cat, error)                      // Cria um novo gato no banco
	GetByID(ctx context.Context, id int64) (domain.Cat, error)                                // Busca um gato pelo ID
	List(ctx context.Context, limit int, cursor *time.Time) ([]domain.Cat, *time.Time, error) // Lista gatos com paginação
	//Update(ctx context.Context, id int64, in domain.CatUpdate) (domain.Cat, error) // Atualiza um gato
	//Delete(ctx context.Context, id int64) error                                   // Remove um gato
}

// CatService define as operações disponíveis para uso externo (ex: API).
type CatService interface {
	Create(ctx context.Context, in domain.CatCreate) (domain.Cat, error)                      // Cria um novo gato
	GetByID(ctx context.Context, id int64) (domain.Cat, error)                                // Busca um gato pelo ID
	List(ctx context.Context, limit int, cursor *time.Time) ([]domain.Cat, *time.Time, error) // Lista gatos
	//Update(ctx context.Context, id int64, in domain.CatUpdate) (domain.Cat, error) // Atualiza um gato
	//Delete(ctx context.Context, id int64) error                                   // Remove um gato
}

// catService é a implementação concreta do CatService.
// Usa um repositório para acessar o banco, um pool de workers para tarefas assíncronas e um timeout para requisições.
type catService struct {
	repo      CatRepository // Repositório para acessar dados dos gatos
	wp        *worker.Pool  // Pool de workers para tarefas assíncronas
	requestTO time.Duration // Tempo limite para cada requisição
}

// NewCatService cria uma nova instância do serviço de gatos.
// Recebe o repositório, o pool de workers e o timeout das requisições.
func NewCatService(repo CatRepository, wp *worker.Pool, requestTimeout time.Duration) CatService {
	return &catService{
		repo:      repo,
		wp:        wp,
		requestTO: requestTimeout,
	}
}

// withTO cria um contexto com timeout para limitar o tempo de execução das operações.
// Se o timeout for zero ou negativo, usa apenas cancelamento manual.
func (s *catService) withTO(ctx context.Context) (context.Context, context.CancelFunc) {
	if s.requestTO <= 0 {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, s.requestTO)
}

// Create cria um novo gato.
// Usa contexto com timeout, chama o repositório para salvar o gato e dispara uma tarefa assíncrona (exemplo: gerar thumbnail).
func (s *catService) Create(ctx context.Context, in domain.CatCreate) (domain.Cat, error) {
	ctx, cancel := s.withTO(ctx)
	defer cancel()

	cat, err := s.repo.Create(ctx, in)
	if err != nil {
		return domain.Cat{}, err
	}

	// Exemplo de tarefa assíncrona: gerar thumbnail (simulado)
	_ = s.wp.Submit(func() error {
		// aqui você faria trabalho pesado (ex.: imagem, chamada externa)
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	return cat, nil
}

// GetByID busca um gato pelo ID.
// Usa contexto com timeout e chama o repositório para buscar o gato.
func (c *catService) GetByID(ctx context.Context, id int64) (domain.Cat, error) {
	ctx, cancel := c.withTO(ctx)
	defer cancel()
	return c.repo.GetByID(ctx, id)
}

// List retorna uma lista de gatos com paginação.
// Usa contexto com timeout e chama o repositório para buscar os gatos.
func (c *catService) List(ctx context.Context, limit int, cursor *time.Time) ([]domain.Cat, *time.Time, error) {
	ctx, cancel := c.withTO(ctx)
	defer cancel()
	return c.repo.List(ctx, limit, cursor)
}

/*
func (s *catService) Update(ctx context.Context, id int64, in domain.UpdateCatInput) (domain.Cat, error) {
    // Atualiza um gato usando contexto com timeout
    ctx, cancel := s.withTO(ctx); defer cancel()
    return s.repo.Update(ctx, id, in)
}

func (s *catService) Delete(ctx context.Context, id int64) error {
    // Remove um gato usando contexto com timeout
    ctx, cancel := s.withTO(ctx); defer cancel()
    return s.repo.Delete(ctx, id)
}
*/

/*
	O defer cancel() serve para garantir que a função cancel() do contexto seja chamada ao final da execução do método, liberando recursos e evitando vazamentos de memória.
	Assim, mesmo se ocorrer erro ou retorno antecipado, o contexto é sempre finalizado corretamente.
*/
