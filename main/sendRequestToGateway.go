package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var httpLock sync.Mutex

func productData(count string, fileName string) []byte {
	s := readFile(fileName)
	times, _ := strconv.Atoi(count)

	var resBuild strings.Builder
	resBuild.WriteString("{\n    \"entityVersion\": \"1.0.0\",\n    \"value\": [")
	// Make different requests with different ID
	for i := 0; i < times; i++ {
		str := strings.ReplaceAll(s, "<%ID%>", strconv.Itoa(i))
		resBuild.WriteString(str + ",")

	}
	res := resBuild.String()
	res = res[:len(res)-1] + "]}"

	return []byte(res)

}

func readFile(fileName string) string {
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println("read fail", err)
	}
	s := string(f)
	return s

}

// OAuth Authentication
func sendHttpRequest(content []byte, token string, namespace string) int {
	// Provide necessary headers
	request, err := http.NewRequest("POST", httpposturl, bytes.NewBuffer(content))
	if err == nil {
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		if namespace != "local" {
			request.Header.Set("Authorization", "Bearer "+token)
		}
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		println(err)
		return 500
	} else {

		log.Println("response Status:", response.StatusCode)
		body, _ := ioutil.ReadAll(response.Body)
		log.Println("response Body:", string(body))
		return response.StatusCode
	}
}

func getToken(namespce string) string {}

type ReqBody struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	Jti         string `json:"jti"`
}

var data []byte

func main() {
	data = productData("10", "conf/10k.json")
	http.HandleFunc("/perf", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(writer http.ResponseWriter, request *http.Request) {
	requests, _ := strconv.Atoi(request.Header.Get("requests"))
	threadNum, _ := strconv.Atoi(request.Header.Get("threadNum"))
	namespace := request.Header.Get("namespace")

	token := getToken(namespace)
	countFail := 0
	var responseTime float64 = 0
	var wg sync.WaitGroup
	limitGoroutine := make(chan int, threadNum)
	println("requests:",requests)

	for i := 0; i < requests; i++ {
		wg.Add(1)
		limitGoroutine <- 0

		go func(i int) {
			t1 := time.Now()
			responseStatus := sendHttpRequest(data, token, namespace)
			t2 := time.Now()
			t3 := t2.Sub(t1).Seconds()
			httpLock.Lock()
			if responseStatus == 202 {
				log.Printf("thread %v :start at %v ,end at %v , cost %v s", i, t1.GoString(), t2.GoString(), t3)
				responseTime += t3
			} else {
				countFail++
			}
			httpLock.Unlock()
			wg.Done()
			<-limitGoroutine
		}(i)
	}
	wg.Wait()
	s := "threadNum:<%threadNum%>, request:<%request%>"
	s = strings.ReplaceAll(s, "<%request%>", strconv.Itoa(requests))
	s = strings.ReplaceAll(s, "<%threadNum%>", strconv.Itoa(threadNum))
	s = strings.ReplaceAll(s, "<%err%>", strconv.Itoa(countFail))
	writer.Write([]byte(s))

}