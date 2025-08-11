package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"load-tester/internal/loadtester"

	"github.com/spf13/cobra"
)

var (
	targetURL   string
	requests    int
	concurrency int
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "load-tester",
		Short: "CLI para testes de carga em serviços web",
		Long: `Load Tester é uma ferramenta CLI para realizar testes de carga em serviços web.
		
Exemplo de uso:
  load-tester --url=http://google.com --requests=1000 --concurrency=10`,
		RunE: runLoadTest,
	}

	// Define flags
	rootCmd.Flags().StringVar(&targetURL, "url", "", "URL do serviço a ser testado (obrigatório)")
	rootCmd.Flags().IntVar(&requests, "requests", 100, "Número total de requests (padrão: 100)")
	rootCmd.Flags().IntVar(&concurrency, "concurrency", 10, "Número de chamadas simultâneas (padrão: 10)")

	// Marca URL como obrigatória
	rootCmd.MarkFlagRequired("url")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runLoadTest(cmd *cobra.Command, args []string) error {
	// Valida parâmetros
	if err := validateParams(); err != nil {
		return err
	}

	// Mostra parâmetros do teste
	printTestParams()

	// Cria configuração
	config := loadtester.Config{
		URL:         targetURL,
		Requests:    requests,
		Concurrency: concurrency,
	}

	// Cria load tester
	tester := loadtester.NewLoadTester(config)

	// Executa teste com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	fmt.Println("Iniciando teste de carga...")
	startTime := time.Now()

	// Executa o teste
	report := tester.Run(ctx)

	fmt.Printf("\nTeste concluído em %v\n", time.Since(startTime))

	// Exibe relatório
	report.PrintReport()

	return nil
}

func validateParams() error {
	// Valida URL
	if targetURL == "" {
		return fmt.Errorf("URL é obrigatória")
	}

	// Verifica se é uma URL válida
	_, err := url.ParseRequestURI(targetURL)
	if err != nil {
		return fmt.Errorf("URL inválida: %v", err)
	}

	// Valida número de requests
	if requests <= 0 {
		return fmt.Errorf("número de requests deve ser maior que 0")
	}

	// Valida concorrência
	if concurrency <= 0 {
		return fmt.Errorf("concorrência deve ser maior que 0")
	}

	// Concorrência não pode ser maior que o número total de requests
	if concurrency > requests {
		concurrency = requests
		fmt.Printf("Aviso: Concorrência ajustada para %d (número total de requests)\n", concurrency)
	}

	return nil
}

func printTestParams() {
	fmt.Println("========== PARÂMETROS DO TESTE ==========")
	fmt.Printf("URL: %s\n", targetURL)
	fmt.Printf("Total de requests: %d\n", requests)
	fmt.Printf("Concorrência: %d\n", concurrency)
	fmt.Println("=========================================")
	fmt.Println()
}
