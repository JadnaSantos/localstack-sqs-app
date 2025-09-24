package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/JadnaSantos/localStack-sqs-app.git/src/internal/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	appCfg := config.Load()
	awsCfg, err := config.AWSConfig(ctx, appCfg)

	if err != nil {
		log.Fatalf("Erro ao configurar AWS SDK: %v", err)
	}

	client := sqs.NewFromConfig(awsCfg)
	queueName := "product-queue"
	result, err := client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})

	if err != nil {
		log.Fatalf("Falha em capturar a url da fila: %v", err)
	}

	for {
		msg, err := client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
			QueueUrl:            result.QueueUrl,
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     5,
			VisibilityTimeout:   10,
		})

		if err != nil {
			log.Fatalf("Falha ao receber mensagem: %v", err)
		}

		if len(msg.Messages) == 0 {
			fmt.Println("Sem mensagens, encerrando programa...")
			break
		}

		for _, message := range msg.Messages {
			log.Printf("Mensagem Recebida: %v \n", *message.Body)
		}
	}
}
