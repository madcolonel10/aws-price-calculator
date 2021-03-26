package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/pricing"
)

func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := pricing.NewFromConfig(cfg)
	desribeServicesOutput, err := client.DescribeServices(context.TODO(), nil)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	fmt.Println(desribeServicesOutput.Services)
	fmt.Println("done")
}
