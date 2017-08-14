package main

import (
	"errors"
	"net/http"
	"time"
)

// Timeout sets the ping timeout in milliseconds
var TimeOut = 30 * time.Second

// PingExec sends a HEAD command to a given URL, returns whether the host answers 200 or not
func PingExec(target string, url string) (httpcode int, err error) {

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err == nil {
		return resp.StatusCode, nil
	}
	//defer resp.Body.Close()
	//fmt.Println(resp.StatusCode)
	return 500, errors.New("Error - Ping when using ping on target " + target)
	//return resp.StatusCode, errors.New(resp.Status)
}
