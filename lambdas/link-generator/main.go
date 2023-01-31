package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"linkshh/pkg/linkshh"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const QueryParamName = "uri"

type ResponseBody struct {
	Uri string `json:"uri"`
}

var logger = log.Default()

func handler(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	logger.Printf("%v", e)

	// get uri url param
	rawUri := e.QueryStringParameters[QueryParamName]
	if rawUri == "" {
		return buildErrorResponse(errors.New("uri is a required parameter"), http.StatusBadRequest), nil
	}

	// convert rawUri to URL to validate
	uri, err := url.Parse(rawUri)
	if err != nil {
		logger.Print(err)
		return buildErrorResponse(errors.New("uri is invalid, please provide a valid uri"), http.StatusBadRequest), nil
	}

	if uri == nil {
		logger.Printf("%s was passed and is invalid", rawUri)
		return buildErrorResponse(errors.New("uri is a required parameter"), http.StatusBadRequest), nil
	}

	shortUri, err := linkshh.ShortenUri(ctx, uri)
	if err != nil {
		logger.Print(err.Error())
		return buildErrorResponse(errors.New("unable to generate url"), http.StatusInternalServerError), nil
	}

	return buildSuccessResponse(shortUri.String()), nil
}

func buildSuccessResponse(uri string) events.APIGatewayProxyResponse {
	b, err := json.Marshal(ResponseBody{Uri: uri})
	if err != nil {
		logger.Printf("error marshalling success response: %s", err.Error())
		return buildErrorResponse(errors.New("unable to generate url"), http.StatusInternalServerError)
	}

	return events.APIGatewayProxyResponse{
		StatusCode:        http.StatusOK,
		Headers:           nil,
		MultiValueHeaders: nil,
		Body:              string(b),
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
