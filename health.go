package main

import "context"

func HealthHandler(ctx context.Context) (*Response, error) {
	return &Response{
		StatusCode: 200,
		Body:       "Ok",
	}, nil
}
