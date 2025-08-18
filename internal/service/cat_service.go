package service

import (
	"context"
	"errors"
	"time"

	"github.com/dya-andrade/cat-api/internal/domain"
	"github.com/dya-andrade/cat-api/internal/worker"
)

var (
	ErrNotFound = errors.New("Cat not found")
)

// CatRepository descreve o que o serviço precisa do repositório.
// (Implementação real está em internal/storage/cat_repo.go)
type CatRepository interface {
	Create(ctx context.Context, in domain.CatCreate) (domain.Cat, error)
	GetByID(ctx context.Context, id int64) (domain.Cat, error)
	List(ctx context.Context, limit int, cursor *time.Time) ([]domain.Cat, *time.Time, error)
	//Update(ctx context.Context, id int64, in domain.CatUpdate) (domain.Cat, error)
	//Delete(ctx context.Context, id int64) error
}

type CatService interface {
	Create(ctx context.Context, in domain.CatCreate) (domain.Cat, error)
	GetByID(ctx context.Context, id int64) (domain.Cat, error)
	List(ctx context.Context, limit int, cursor *time.Time) ([]domain.Cat, *time.Time, error)
	//Update(ctx context.Context, id int64, in domain.CatUpdate) (domain.Cat, error)
	//Delete(ctx context.Context, id int64) error
}

type catService struct {
	repo      CatRepository
	wp        *worker.Pool
	requestTO time.Duration
}

func NewCatService(repo CatRepository, wp *worker.Pool, requestTimeout time.Duration) CatService {
	return &catService{
		repo:      repo,
		wp:        wp,
		requestTO: requestTimeout,
	}
}
