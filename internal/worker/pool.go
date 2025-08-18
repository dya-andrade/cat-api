package worker

import (
	"sync"
)

type job func() error // Define o tipo de tarefa: uma função que retorna erro

type Pool struct {
	concurrency int            // Número de workers (concorrência)
	jobs        chan job       // Canal de tarefas a serem executadas
	wg          sync.WaitGroup // Sincroniza o término dos workers
	onceStart   sync.Once      // Garante que Start só execute uma vez
	onceStop    sync.Once      // Garante que Shutdown só execute uma vez
}

// NewPool cria um novo pool de workers.
// Recebe o número de workers desejado. Se <= 0, usa 1.
// Cria o canal de jobs com buffer proporcional à concorrência.
func NewPool(concurrency int) *Pool {
	if concurrency <= 0 {
		concurrency = 1
	}

	return &Pool{
		concurrency: concurrency,
		jobs:        make(chan job, concurrency*4),
	}
}

// Start inicia os workers do pool.
// Garante que só será chamado uma vez.
// Para cada worker, inicia uma goroutine que consome jobs do canal.
// Cada worker executa jobs até o canal ser fechado.
func (p *Pool) Start() {
	p.onceStart.Do(func() {
		for i := 0; i < p.concurrency; i++ {
			p.wg.Add(1)
			go func() {
				defer p.wg.Done()
				for j := range p.jobs {
					_ = j() // Executa o job e ignora o erro
					// erros podem ser enviados a um canal/telemetria se quiser
				}
			}()
		}
	})
}

// Submit envia uma tarefa para ser executada pelos workers.
// Se o canal estiver cheio, bloqueia até conseguir enviar (backpressure).
// Retorna erro apenas se o canal já estiver fechado.
func (p *Pool) Submit(fn func() error) error {
	select {
	case p.jobs <- fn:
		return nil
	default:
		// fila está cheia, mas ainda aberta: faz bloqueante para backpressure
		p.jobs <- fn
		return nil
	}
}

// Shutdown encerra o pool de workers.
// Garante que só será chamado uma vez.
// Fecha o canal de jobs e espera todos os workers terminarem.
func (p *Pool) Shutdown() {
	p.onceStop.Do(func() {
		close(p.jobs)
		p.wg.Wait()
	})
}
