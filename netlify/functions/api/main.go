package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const targetURL = "https://chipper-biscotti-c69d96.netlify.app/.netlify/functions/api"
const concurrentRequests = 100

func makeRequest(wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(targetURL)
	if err != nil {
		log.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("Request successful")
	} else {
		log.Printf("Request returned status code: %d\n", resp.StatusCode)
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	log.Println("Starting concurrent requests...")

	var wg sync.WaitGroup
	wg.Add(concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go makeRequest(&wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("All requests completed.")

	response := &events.APIGatewayProxyResponse{
		StatusCode:      200,
		Headers:         map[string]string{"Content-Type": "text/plain"},
		Body:            "Successfully completed concurrent requests.",
		IsBase64Encoded: false,
	}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
