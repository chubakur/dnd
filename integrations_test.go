package main

import (
	"io"
	"net/http"
	"os"
	"testing"
)

func TestLocalHealth(t *testing.T) {

}

func TestProdHealth(t *testing.T) {
	healthUrl := os.Getenv("HEALTH_CHECK_PROD_URL")
	resp, err := http.Get(healthUrl)
	if err != nil {
		t.Fatal(err)
	}
	bRes, err := io.ReadAll(resp.Body)
	strRes := string(bRes)
	if err != nil {
		t.Error(err)
	}
	if strRes != "Ok" {
		t.Errorf("Integration health-check error. Ok != %s", strRes)
	}
}
