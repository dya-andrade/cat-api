package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dya-andrade/cat-api/internal/http/handlers"
	"github.com/dya-andrade/cat-api/internal/service"
)

func NewRouter(catSvc service.CatService) http.Handler {
	r := chi.NewRouter() // Cria um novo roteador usando o chi

	// middlewares essenciais
	r.Use(middleware.RequestID)          // Adiciona um ID único para cada requisição (útil para rastreamento)
	r.Use(middleware.RealIP)             // Captura o IP real do cliente (mesmo atrás de proxy)
	r.Use(middleware.Logger)             // Loga informações de cada requisição (método, rota, tempo, etc)
	r.Use(middleware.Recoverer)          // Recupera de panics e retorna erro 500 ao invés de travar o servidor
	r.Use(middleware.Heartbeat("/live")) // Endpoint simples para checagem de vida (/live)

	// health/ready
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)              // Responde com status 200 OK
		_, _ = w.Write([]byte(`{"status":"ok"}`)) // Retorna JSON simples indicando que está saudável
	})

	// handlers
	cats := handlers.NewCatsHandler(catSvc) // Cria o handler dos gatos, passando o serviço

	r.Route("/cats", func(r chi.Router) {
		r.Get("/", cats.List)        // GET /cats?limit=...&cursor=RFC3339 -> lista gatos
		r.Post("/", cats.Create)     // POST /cats -> cria novo gato
		r.Get("/{id}", cats.GetByID) // GET /cats/{id} -> busca gato por ID
		//r.Put("/{id}", cats.Update)    // PUT /cats/{id} -> atualiza gato (comentado)
		//r.Delete("/{id}", cats.Delete) // DELETE /cats/{id} -> remove gato (comentado)
	})

	return r // Retorna o roteador configurado
}

/*
	O que são middlewares?
	Middlewares são funções que interceptam e processam requisições HTTP antes ou depois dos handlers principais.
	Eles podem adicionar funcionalidades como logs, autenticação, tratamento de erros, rastreamento, etc.,
	sem alterar o código dos handlers.
*/
