# Load Tester CLI - Documentação

## Como Funciona

O Load Tester é uma ferramenta CLI que realiza testes de carga em serviços web através de requisições HTTP GET concorrentes.

### Algoritmo de Funcionamento

1. **Validação de Parâmetros**: Verifica URL, número de requests e concorrência
2. **Controle de Concorrência**: Usa semáforos para limitar requisições simultâneas
3. **Execução Paralela**: Distribui requests entre goroutines conforme concorrência
4. **Coleta de Métricas**: Registra tempo de resposta e status code de cada request
5. **Geração de Relatório**: Compila estatísticas e apresenta resultados

### Controle de Concorrência

```go
// Semáforo limita requisições simultâneas
semaphore := make(chan struct{}, concurrency)

// Cada goroutine adquire o semáforo
semaphore <- struct{}{}        // Adquire
defer func() { <-semaphore }() // Libera
```

## Parâmetros de Entrada

### Flags CLI

| Flag | Descrição | Obrigatório | Padrão |
|------|-----------|-------------|---------|
| `--url` | URL do serviço a ser testado | ✅ | - |
| `--requests` | Número total de requisições | ❌ | 100 |
| `--concurrency` | Chamadas simultâneas | ❌ | 10 |

### Validações

- **URL**: Deve ser uma URL válida (http/https)
- **Requests**: Maior que 0
- **Concurrency**: Maior que 0 e não superior ao total de requests

## Execução

### Exemplo Básico
```bash
./load-tester --url=http://google.com --requests=1000 --concurrency=10
```

### Com Docker
```bash
docker run <imagem> --url=http://google.com --requests=1000 --concurrency=10
```

### Cenários de Teste

**Teste de Stress:**
```bash
--url=https://api.exemplo.com --requests=5000 --concurrency=100
```

**Teste de Estabilidade:**
```bash
--url=https://api.exemplo.com --requests=1000 --concurrency=5
```

**Teste de Latência:**
```bash
--url=https://api.exemplo.com --requests=100 --concurrency=1
```

## Relatório Gerado

### Métricas Coletadas

**Tempo de Execução:**
- Tempo total do teste
- Tempo médio de resposta por request
- Tempo mínimo e máximo de resposta

**Contadores de Status:**
- Total de requisições realizadas
- Requisições com status 200 (sucesso)
- Distribuição completa de códigos de status
- Contagem de erros de conexão/timeout

**Performance:**
- Taxa de sucesso (%)
- Requisições por segundo
- Erros de rede/timeout

### Exemplo de Relatório

```
========== RELATÓRIO DE TESTE DE CARGA ==========
Tempo total de execução: 2.345s
Total de requisições: 1000
Requisições com status 200: 950
Taxa de sucesso: 95.00%

Distribuição de códigos de status:
  Status 200: 950 requisições
  Status 404: 30 requisições
  Status 500: 20 requisições

Estatísticas de tempo de resposta:
  Tempo médio: 45ms
  Tempo mínimo: 12ms
  Tempo máximo: 890ms
  Requisições por segundo: 426.50

================================================
```

## Configurações Avançadas

### Timeout de Requisição

Por padrão, cada requisição tem timeout de 30 segundos:

```go
client: &http.Client{
    Timeout: 30 * time.Second,
}
```

### Timeout do Teste Completo

O teste completo tem timeout de 10 minutos:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
```

### Headers Customizados

Atualmente usa apenas GET requests. Para headers customizados, seria necessário estender a configuração:

```go
req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
// req.Header.Set("Authorization", "Bearer token")
```

## Casos de Uso

### 1. Teste de Capacidade
Determina quantas requisições simultâneas o serviço suporta:

```bash
# Incrementa concorrência até encontrar limite
./load-tester --url=https://api.exemplo.com --requests=1000 --concurrency=50
./load-tester --url=https://api.exemplo.com --requests=1000 --concurrency=100
./load-tester --url=https://api.exemplo.com --requests=1000 --concurrency=200
```

### 2. Teste de Regressão
Compara performance entre versões:

```bash
# Antes do deploy
./load-tester --url=https://staging.exemplo.com --requests=500 --concurrency=25

# Após o deploy  
./load-tester --url=https://production.exemplo.com --requests=500 --concurrency=25
```

### 3. Monitoramento de SLA
Verifica se API atende requisitos de latência:

```bash
# Deve ter tempo médio < 100ms com 95% de sucesso
./load-tester --url=https://api.exemplo.com --requests=200 --concurrency=10
```

## Interpretação dos Resultados

### Taxa de Sucesso
- **> 99%**: Excelente
- **95-99%**: Bom  
- **90-95%**: Aceitável
- **< 90%**: Problemático

### Tempo de Resposta
- **< 100ms**: Muito rápido
- **100-500ms**: Rápido
- **500ms-2s**: Aceitável
- **> 2s**: Lento

### Requisições por Segundo
- Métrica de throughput
- Compara com requisitos de capacidade
- Identifica gargalos de performance
