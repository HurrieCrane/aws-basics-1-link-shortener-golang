package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"linkshh/pkg/linkshh"
	"log"
	"net/http"
	"net/url"
)

const QueryParamName = "hash"

var logger = log.Default()

func handler(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.Printf("%v", e)

	// get hash url param
	hash := e.PathParameters[QueryParamName]
	if hash == "" {
		return buildErrorResponse(errors.New("uri is a required parameter"), http.StatusBadRequest), nil
	}

	tinyUri, err := linkshh.ExpandHash(ctx, hash)
	if err != nil {
		logger.Print(err.Error())
		return buildErrorResponse(errors.New("unable to generate url"), http.StatusInternalServerError), nil
	}
	return buildSuccessResponse(tinyUri), nil
}

func buildSuccessResponse(tinyUri *url.URL) events.APIGatewayProxyResponse {
	logger.Printf("url to return: %s", tinyUri.String())
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMovedPermanently,
		Headers: map[string]string{
			"Location": tinyUri.String(),
		},
		MultiValueHeaders: nil,
		Body:              "",
		IsBase64Encoded:   false,
	}
}

func buildErrorResponse(e error, statusCode int) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:        statusCode,
		Headers:           nil,
		MultiValueHeaders: nil,
		Body:              fmt.Sprintf("{ \"errorMsg\": \"%s\" }", e.Error()),
		IsBase64Encoded:   false,
	}
}

func main() {
	lambda.Start(handler)
}
