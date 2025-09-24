package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/JadnaSantos/localStack-sqs-app.git/src/internal/config"
	resterr "github.com/JadnaSantos/localStack-sqs-app.git/src/internal/rest_err"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/labstack/echo/v4"
)

type Product struct {
	Price float64 `json:"price"`
	Name  string  `json:"name"`
}

func productHandle(c echo.Context) error {
	var p Product

	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, resterr.NewBadRequestError("Bad Request"))
	}

	if err := publishInQueue(p); err != nil {
		return c.JSON(http.StatusInternalServerError, resterr.NewInternalServerError("Internal Server Error"))
	}

	return c.NoContent(202)
}

func publishInQueue(p Product) error {
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

	payload, err := json.Marshal(p)

	if err != nil {
		log.Fatalf("falha ao parsear struct: %v", err)
	}

	out, err := client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    result.QueueUrl,
		MessageBody: aws.String(string(payload)),
	})

	if err != nil {
		log.Fatalf("falha ao enviar mensagem: %v", err)
	}

	fmt.Println(*out.MessageId)

	return nil
}

func main() {
	e := echo.New()
	e.POST("/product", productHandle)
	e.Logger.Fatal(e.Start(":3000"))
}
