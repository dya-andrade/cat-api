package storage

import (
	"context"
	"errors"
	"time"

	"github.com/dya-andrade/cat-api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
* indica ponteiro, ou seja, uma referência para o valor de uma variável (pode ser nulo).
& obtém o endereço de memória de uma variável, criando um ponteiro para ela.
Exemplo:

*CatRepository é um ponteiro para um objeto CatRepository.
&CatRepository{db: db} cria um ponteiro para a struct inicializada.
*/

type CatRepository struct {
	db *pgxpool.Pool // Conexão com o banco de dados PostgreSQL
}

// Cria uma nova instância de CatRepository usando o pool de conexões
func NewCatRepository(db *pgxpool.Pool) *CatRepository {
	return &CatRepository{db: db} // Retorna o repositório com o banco configurado
}

/*
O contexto (context.Context) é necessário para controlar o tempo de execução, cancelamento e deadlines de operações no banco de dados.
Ele permite, por exemplo, cancelar uma consulta se ela demorar demais ou se a requisição do usuário for encerrada.
Assim, evita travamentos e libera recursos corretamente.
*/

func (repository *CatRepository) Create(ctx context.Context, in domain.CatCreate) (domain.Cat, error) {
	// Cria um novo registro de gato no banco de dados

	row := repository.db.QueryRow(
		ctx,
		// Executa o comando SQL para inserir um novo gato e retorna os dados inseridos
		"INSERT INTO cats (name, age_years, breed, coat_color, weight_kg) VALUES ($1,$2,$3,$4,$5) RETURNING id, name, age_years, breed, coat_color, weight_kg, created_at, updated_at",
		in.Name, in.AgeYears, in.Breed, in.CoatColor, in.WeightKG,
		// Passa os valores do novo gato para os parâmetros da query
	)

	var c domain.Cat
	// Cria uma variável para armazenar o resultado retornado do banco

	err := row.Scan(&c.ID, &c.Name, &c.AgeYears, &c.Breed, &c.CoatColor, &c.WeightKG, &c.CreatedAt, &c.UpdatedAt)
	// Lê os dados retornados pela query e preenche a struct Cat

	return c, err
	// Retorna o gato criado e um erro (se houver)
}

func (repository *CatRepository) GetByID(ctx context.Context, id int64) (domain.Cat, error) {
	// Busca um gato pelo ID no banco de dados

	row := repository.db.QueryRow(
		ctx,
		// Executa o comando SQL para selecionar o gato pelo ID
		"SELECT id, name, age_years, breed, coat_color, weight_kg, created_at, updated_at FROM cats WHERE id=$1",
		id,
		// Passa o ID como parâmetro para a query
	)

	var c domain.Cat
	// Cria uma variável para armazenar o resultado retornado do banco

	if err := row.Scan(&c.ID, &c.Name, &c.AgeYears, &c.Breed, &c.CoatColor, &c.WeightKG, &c.CreatedAt, &c.UpdatedAt); err != nil {
		// Lê os dados retornados pela query e preenche a struct Cat
		if errors.Is(err, pgx.ErrNoRows) {
			// Se não encontrar nenhum registro, retorna struct vazia e erro
			return domain.Cat{}, err
		}
		// Se ocorrer outro erro, retorna struct vazia e erro
		return domain.Cat{}, err
	}

	return c, nil
	// Retorna o gato encontrado e nil para erro
}

// Paginação por cursor baseado em created_at/id
func (repository *CatRepository) List(ctx context.Context, limit int, cursor *time.Time) ([]domain.Cat, *time.Time, error) {
	// Lista gatos com paginação usando cursor

	var rows pgx.Rows
	var err error

	if cursor == nil {
		// Se não houver cursor, busca os primeiros registros
		rows, err = repository.db.Query(
			ctx,
			"SELECT id, name, age_years, breed, coat_color, weight_kg, created_at, updated_at FROM cats ORDER BY created_at DESC LIMIT $1",
			limit,
		)
	} else {
		// Se houver cursor, busca registros após o cursor
		rows, err = repository.db.Query(
			ctx,
			"SELECT id, name, age_years, breed, coat_color, weight_kg, created_at, updated_at FROM cats WHERE created_at < $1 ORDER BY created_at DESC LIMIT $2",
			*cursor,
			limit,
		)
	}

	// Verifica se houve erro na consulta
	if err != nil {
		return nil, nil, err
	}

	// Fecha as linhas após a consulta
	defer rows.Close()

	var cats []domain.Cat
	var lastCreatedAt *time.Time

	for rows.Next() {
		var c domain.Cat
		if err := rows.Scan(&c.ID, &c.Name, &c.AgeYears, &c.Breed, &c.CoatColor, &c.WeightKG, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, nil, err
		}
		cats = append(cats, c)
		lastCreatedAt = &c.CreatedAt // Atualiza o último created_at encontrado
	}

	if len(cats) == 0 {
		return nil, nil, nil // Retorna nil se não encontrar gatos
	}

	return cats, lastCreatedAt, nil // Retorna a lista de gatos e o último created_at
}
