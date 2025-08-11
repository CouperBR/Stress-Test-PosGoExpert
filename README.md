# Load Tester CLI

Sistema CLI para realizar testes de carga em serviços web.

## Como Usar

### Parâmetros

| Flag | Descrição | Obrigatório | Padrão |
|------|-----------|-------------|---------|
| `--url` | URL do serviço a ser testado | ✅ | - |
| `--requests` | Número total de requisições | ❌ | 100 |
| `--concurrency` | Chamadas simultâneas | ❌ | 10 |

### Executando Localmente

```bash
# Compila
go build -o load-tester cmd/main.go

# Executa
./load-tester --url=https://google.com --requests=1000 --concurrency=10
```

### Executando com Docker

```bash
# Constrói imagem
docker build -t load-tester .

# Executa teste
docker run load-tester --url=http://google.com --requests=1000 --concurrency=10
```

## Exemplos

```bash
# Teste básico
./load-tester --url=https://httpbin.org/status/200 --requests=50 --concurrency=5

# Teste de alta carga
./load-tester --url=https://api.exemplo.com --requests=5000 --concurrency=100

# Teste de latência
./load-tester --url=https://site.com --requests=100 --concurrency=1
```

## Relatório

O sistema gera um relatório com:
- Tempo total de execução
- Total de requisições realizadas
- Requisições com status 200
- Distribuição de códigos de status
- Estatísticas de tempo de resposta
