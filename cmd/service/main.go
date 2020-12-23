package main

import (
	"errors"
	"fmt"
	"gopkg.in/robfig/cron.v2"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

func main() {
	isServerWorking := validateIfServerIsListening()
	if !isServerWorking {
		panic(errors.New("server is not listening"))
	}

	c := cron.New()
	c.AddFunc("@hourly", callExternalResource)

	go c.Start()
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func callExternalResource() {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/users/fallback", nil)
	if err != nil {
		fmt.Println("error constructing the request: " + err.Error())
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error executing the fallback method: " + err.Error())
		panic(err)
	}
	if resp.StatusCode != 200 {
		fmt.Println("error during the call, the server returned: " + strconv.Itoa(resp.StatusCode))
		panic(errors.New("the server returned: "  + strconv.Itoa(resp.StatusCode)))
	}
}

func validateIfServerIsListening() bool {
	resp, err := http.Get("http://localhost:8080/api/health")
	if err != nil {
		fmt.Println("error getting the health: " + err.Error())
		return false
	}

	if resp.StatusCode != 200 {
		fmt.Println("error: the response of the server is " + strconv.Itoa(resp.StatusCode))
		return false
	}

	return true
}
