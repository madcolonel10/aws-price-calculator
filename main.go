package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/service/pricing"
// 	"github.com/aws/aws-sdk-go-v2/service/pricing/types"
// 	"github.com/aws/aws-sdk-go/aws"
// )

// func main() {
// 	// Load the Shared AWS Configuration (~/.aws/config)
// 	cfg, err := config.LoadDefaultConfig(context.TODO())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	client := pricing.NewFromConfig(cfg)
// 	serviceCode := "AmazonEC2"
// 	desribeServicesOutput, err := client.DescribeServices(context.TODO(), &pricing.DescribeServicesInput{ServiceCode: &serviceCode})
// 	if err != nil {
// 		fmt.Printf("error: %s\n", err)
// 		return
// 	}
// 	fmt.Println(desribeServicesOutput.Services[0].AttributeNames)
// 	fmt.Println("done")
// 	//cpu ram storage instanceType operating-system (rhel, centos, ubuntu)
// 	input := &pricing.GetProductsInput{
// 		Filters: []types.Filter{
// 			{
// 				Field: aws.String("ServiceCode"),
// 				Type:  types.FilterTypeTermMatch,
// 				Value: aws.String("AmazonEC2"),
// 			},
// 			{
// 				Field: aws.String("instanceType"),
// 				Type:  types.FilterTypeTermMatch,
// 				Value: aws.String("t2.medium"),
// 			},
// 			{
// 				Field: aws.String("vcpu"),
// 				Type:  types.FilterTypeTermMatch,
// 				Value: aws.String("2"),
// 			},
// 			{
// 				Field: aws.String("memory"),
// 				Type:  types.FilterTypeTermMatch,
// 				Value: aws.String("4 GiB"),
// 			},
// 		},
// 		FormatVersion: aws.String("aws_v1"),
// 		MaxResults:    5,
// 		ServiceCode:   aws.String("AmazonEC2"),
// 	}

// 	result, err := client.GetProducts(context.TODO(), input)
// 	if err != nil {
// 		fmt.Println("error here")
// 		fmt.Println(err.Error())
// 		return
// 	}

// 	jsonResult, err := json.Marshal(result)

// 	fmt.Println(string(jsonResult))
// }
