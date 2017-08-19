package main_test

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

// TestIntegration - only check localhost:port#
func TestIntegration(t *testing.T) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	resp, err := http.Get("https://httpbin.org/status/418")
	if err != nil || resp.StatusCode != 418 {
		fmt.Printf("Failed with status 418 resp %v error %v \n", resp, err)
		t.Fail()
	}
	//
	//resp1, err1 := http.Get("http://localhost:9000/pingqa")
	//if err1 != nil || resp1.StatusCode != 200 {
	//	fmt.Printf("Failed with localhost resp %v error %v \n", resp1, err1)
	//	t.Fail()
	//}
}
