package loadtester

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoadTester_SingleRequest(t *testing.T) {
	// Cria servidor de teste
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Configura teste
	config := Config{
		URL:         server.URL,
		Requests:    1,
		Concurrency: 1,
	}

	tester := NewLoadTester(config)
	ctx := context.Background()

	// Executa teste
	report := tester.Run(ctx)

	// Verifica resultados
	if report.TotalRequests != 1 {
		t.Errorf("Esperado 1 request, obtido %d", report.TotalRequests)
	}

	if report.SuccessfulReqs != 1 {
		t.Errorf("Esperado 1 request bem-sucedido, obtido %d", report.SuccessfulReqs)
	}

	if report.StatusCodeCounts[200] != 1 {
		t.Errorf("Esperado 1 status 200, obtido %d", report.StatusCodeCounts[200])
	}

	if report.ErrorCount != 0 {
		t.Errorf("Esperado 0 erros, obtido %d", report.ErrorCount)
	}
}

func TestLoadTester_MultipleRequests(t *testing.T) {
	// Cria servidor de teste
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Configura teste
	config := Config{
		URL:         server.URL,
		Requests:    10,
		Concurrency: 3,
	}

	tester := NewLoadTester(config)
	ctx := context.Background()

	// Executa teste
	report := tester.Run(ctx)

	// Verifica resultados
	if report.TotalRequests != 10 {
		t.Errorf("Esperado 10 requests, obtido %d", report.TotalRequests)
	}

	if report.SuccessfulReqs != 10 {
		t.Errorf("Esperado 10 requests bem-sucedidos, obtido %d", report.SuccessfulReqs)
	}

	if report.StatusCodeCounts[200] != 10 {
		t.Errorf("Esperado 10 status 200, obtido %d", report.StatusCodeCounts[200])
	}
}

func TestLoadTester_DifferentStatusCodes(t *testing.T) {
	requestCount := 0

	// Cria servidor que alterna entre status codes
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		if requestCount%2 == 0 {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		w.Write([]byte("Response"))
	}))
	defer server.Close()

	// Configura teste
	config := Config{
		URL:         server.URL,
		Requests:    6,
		Concurrency: 2,
	}

	tester := NewLoadTester(config)
	ctx := context.Background()

	// Executa teste
	report := tester.Run(ctx)

	// Verifica resultados
	if report.TotalRequests != 6 {
		t.Errorf("Esperado 6 requests, obtido %d", report.TotalRequests)
	}

	// Deve ter tanto 200 quanto 404
	if report.StatusCodeCounts[200] == 0 {
		t.Error("Esperado pelo menos um status 200")
	}

	if report.StatusCodeCounts[404] == 0 {
		t.Error("Esperado pelo menos um status 404")
	}
}

func TestLoadTester_WithDelay(t *testing.T) {
	// Cria servidor com delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Configura teste
	config := Config{
		URL:         server.URL,
		Requests:    3,
		Concurrency: 1,
	}

	tester := NewLoadTester(config)
	ctx := context.Background()

	start := time.Now()
	report := tester.Run(ctx)
	elapsed := time.Since(start)

	// Verifica que levou pelo menos 300ms (3 requests x 100ms cada)
	if elapsed < 300*time.Millisecond {
		t.Errorf("Teste muito rápido, esperado pelo menos 300ms, obtido %v", elapsed)
	}

	// Verifica estatísticas de tempo
	if report.AvgResponseTime < 90*time.Millisecond {
		t.Errorf("Tempo médio muito baixo, esperado ~100ms, obtido %v", report.AvgResponseTime)
	}
}

func TestLoadTester_Concurrency(t *testing.T) {
	// Cria servidor com delay
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Teste com concorrência alta
	config := Config{
		URL:         server.URL,
		Requests:    10,
		Concurrency: 10, // Todas ao mesmo tempo
	}

	tester := NewLoadTester(config)
	ctx := context.Background()

	start := time.Now()
	report := tester.Run(ctx)
	elapsed := time.Since(start)

	// Com concorrência 10, deve levar cerca de 200ms (não 2000ms)
	if elapsed > 1*time.Second {
		t.Errorf("Concorrência não funcionou, muito lento: %v", elapsed)
	}

	if report.TotalRequests != 10 {
		t.Errorf("Esperado 10 requests, obtido %d", report.TotalRequests)
	}
}

func TestReport_PrintReport(t *testing.T) {
	// Testa se PrintReport não quebra
	report := &Report{
		TotalTime:        1 * time.Second,
		TotalRequests:    100,
		SuccessfulReqs:   95,
		StatusCodeCounts: map[int]int{200: 95, 404: 3, 500: 2},
		ErrorCount:       0,
		AvgResponseTime:  50 * time.Millisecond,
		MinResponseTime:  10 * time.Millisecond,
		MaxResponseTime:  200 * time.Millisecond,
	}

	// Não deve quebrar
	report.PrintReport()
}
