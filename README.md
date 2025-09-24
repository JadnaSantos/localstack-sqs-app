# localStack-sqs-app

Aplicativo de demonstração que publica e consome mensagens SQS usando LocalStack. Ele expõe um produtor HTTP que recebe pedidos de criação de produtos e um consumidor CLI que lê essas mensagens da fila `product-queue`.

## Requisitos

- Go 1.21 ou superior
- Docker e Docker Compose
- AWS CLI ou awslocal (opcional, para criação/inspeção da fila)

## Configuração

1. Carregue essas variáveis no terminal que executará os binários:
   ```bash
   export $(grep -v '^#' .env | xargs)
   ```
2. Suba o LocalStack com Docker Compose:
   ```bash
   docker compose up -d
   ```
3. Crie a fila SQS dentro do LocalStack (via awslocal ou AWS CLI):
   ```bash
   awslocal sqs create-queue --queue-name product-queue
   # ou
   aws --endpoint-url http://localhost:port sqs create-queue --queue-name product-queue
   ```

## Produtor HTTP

O produtor sobe um servidor Echo na porta `3000` e publica o corpo do POST na fila.

```bash
go run ./src/cmd/producer/main.go
```

Exemplo de requisição:

```bash
curl -X POST http://localhost:3000/product \
  -H 'Content-Type: application/json' \
  -d '{"name":"Keyboard","price":199.9}'
```

Uma resposta `202 Accepted` indica que a mensagem foi enfileirada.

## Consumidor CLI

O consumidor busca mensagens da mesma fila e imprime o corpo recebido.

```bash
go run ./src/cmd/consumer/main.go
```

Quando a fila estiver vazia, o processo encerra com a mensagem `Sem mensagens, encerrando programa...`.

## Utilidades adicionais

- Listar filas disponíveis:
  ```bash
  awslocal sqs list-queues
  ```
- Visualizar mensagens pendentes:
  ```bash
  awslocal sqs receive-message --queue-url http://localhost:4566/000000000000/product-queue
  ```
- Derrubar o LocalStack:
  ```bash
  docker compose down
  ```

## Desenvolvimento

- Atualize os módulos quando necessário:
  ```bash
  go mod tidy
  ```
- Rodar o build completo:
  ```bash
  go build ./...
  ```
