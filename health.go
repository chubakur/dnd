package main

import "context"

func HealthHandler(ctx context.Context) *Response {
	return &Response{
		StatusCode: 200,
		Body:       "Ok",
	}
}
