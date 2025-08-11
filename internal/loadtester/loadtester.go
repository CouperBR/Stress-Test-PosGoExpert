package loadtester

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Config contém os parâmetros do teste de carga
type Config struct {
	URL         string
	Requests    int
	Concurrency int
}

// Result armazena os resultados de uma requisição
type Result struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

// Report contém o relatório final do teste
type Report struct {
	TotalTime        time.Duration
	TotalRequests    int
	SuccessfulReqs   int
	StatusCodeCounts map[int]int
	ErrorCount       int
	AvgResponseTime  time.Duration
	MinResponseTime  time.Duration
	MaxResponseTime  time.Duration
}

// LoadTester executa testes de carga
type LoadTester struct {
	config Config
	client *http.Client
}

// NewLoadTester cria uma nova instância do load tester
func NewLoadTester(config Config) *LoadTester {
	return &LoadTester{
		config: config,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Run executa o teste de carga
func (lt *LoadTester) Run(ctx context.Context) *Report {
	startTime := time.Now()

	// Canal para resultados
	results := make(chan Result, lt.config.Requests)

	// Canal para controlar concorrência
	semaphore := make(chan struct{}, lt.config.Concurrency)

	var wg sync.WaitGroup

	// Lança as requisições
	for i := 0; i < lt.config.Requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Adquire semáforo
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Executa requisição
			result := lt.makeRequest(ctx)
			results <- result
		}()
	}

	// Goroutine para fechar canal de resultados
	go func() {
		wg.Wait()
		close(results)
	}()

	// Coleta resultados
	return lt.collectResults(results, startTime)
}

// makeRequest executa uma única requisição HTTP
func (lt *LoadTester) makeRequest(ctx context.Context) Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", lt.config.URL, nil)
	if err != nil {
		return Result{
			Error:    err,
			Duration: time.Since(start),
		}
	}

	resp, err := lt.client.Do(req)
	if err != nil {
		return Result{
			Error:    err,
			Duration: time.Since(start),
		}
	}
	defer resp.Body.Close()

	return Result{
		StatusCode: resp.StatusCode,
		Duration:   time.Since(start),
	}
}

// collectResults coleta e processa todos os resultados
func (lt *LoadTester) collectResults(results <-chan Result, startTime time.Time) *Report {
	report := &Report{
		StatusCodeCounts: make(map[int]int),
		MinResponseTime:  time.Hour, // Valor inicial alto
	}

	var totalDuration time.Duration

	for result := range results {
		report.TotalRequests++

		if result.Error != nil {
			report.ErrorCount++
		} else {
			report.StatusCodeCounts[result.StatusCode]++
			if result.StatusCode == 200 {
				report.SuccessfulReqs++
			}
		}

		// Estatísticas de tempo de resposta
		totalDuration += result.Duration
		if result.Duration < report.MinResponseTime {
			report.MinResponseTime = result.Duration
		}
		if result.Duration > report.MaxResponseTime {
			report.MaxResponseTime = result.Duration
		}
	}

	report.TotalTime = time.Since(startTime)

	// Calcula tempo médio de resposta
	if report.TotalRequests > 0 {
		report.AvgResponseTime = totalDuration / time.Duration(report.TotalRequests)
	}

	// Se não houve requisições válidas, zera o tempo mínimo
	if report.MinResponseTime == time.Hour {
		report.MinResponseTime = 0
	}

	return report
}

// PrintReport imprime o relatório formatado
func (r *Report) PrintReport() {
	fmt.Println("\n========== RELATÓRIO DE TESTE DE CARGA ==========")
	fmt.Printf("Tempo total de execução: %v\n", r.TotalTime)
	fmt.Printf("Total de requisições: %d\n", r.TotalRequests)
	fmt.Printf("Requisições com status 200: %d\n", r.SuccessfulReqs)
	fmt.Printf("Taxa de sucesso: %.2f%%\n", float64(r.SuccessfulReqs)/float64(r.TotalRequests)*100)

	if r.ErrorCount > 0 {
		fmt.Printf("Erros de conexão/timeout: %d\n", r.ErrorCount)
	}

	fmt.Println("\nDistribuição de códigos de status:")
	for statusCode, count := range r.StatusCodeCounts {
		fmt.Printf("  Status %d: %d requisições\n", statusCode, count)
	}

	fmt.Println("\nEstatísticas de tempo de resposta:")
	fmt.Printf("  Tempo médio: %v\n", r.AvgResponseTime)
	fmt.Printf("  Tempo mínimo: %v\n", r.MinResponseTime)
	fmt.Printf("  Tempo máximo: %v\n", r.MaxResponseTime)

	// Calcula requisições por segundo
	reqsPerSec := float64(r.TotalRequests) / r.TotalTime.Seconds()
	fmt.Printf("  Requisições por segundo: %.2f\n", reqsPerSec)

	fmt.Println("================================================")
}
