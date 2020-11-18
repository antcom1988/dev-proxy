package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var ports map[string]string

func main() {
	readPorts()

	serverMuxA := http.NewServeMux()
	serverMuxA.HandleFunc("/", defaultHandler)

	serverMuxB := http.NewServeMux()
	serverMuxB.HandleFunc("/", handler)

	fmt.Println("Starting Proxy Service ...")

	for key, val := range ports {
		fmt.Println("port " + key + " redirect to "+ val)
		go http.ListenAndServe(":"+key, serverMuxB)
	}

	http.ListenAndServe(":1000", serverMuxA)
}

func readPorts(){
	file, err := os.Open("ports.txt")
	fatal(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	ports = map[string]string{}
	for scanner.Scan() {
		mapPort := strings.Split(scanner.Text(), ";")
		ports[mapPort[0]] = mapPort[1]
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request){
	localAddr :=  r.Context().Value(http.LocalAddrContextKey).(*net.TCPAddr)
	port := strconv.Itoa(localAddr.Port)
	targetUrl := ports[port]

	fmt.Println("HOST : " + r.Host)
	fmt.Println("PORT : " + port)
	fmt.Println("URL : " + targetUrl)

	uri := targetUrl+r.RequestURI
	fmt.Println(r.Method + ": " + uri)

	var jsonStr = []byte("")

	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		body, err := ioutil.ReadAll(r.Body)
		fatal(err)
		fmt.Printf("Body: %v\n", string(body))

		jsonStr = []byte(body)
	}

	rr, err := http.NewRequest(r.Method, uri, bytes.NewBuffer(jsonStr))
	fatal(err)

	copyHeader(r.Header, &rr.Header)

	var transport http.Transport
	resp, err := transport.RoundTrip(rr)
	fatal(err)

	fmt.Printf("Resp-Headers: %v\n", resp.Header)
	fmt.Println("====================================================================")

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fatal(err)

	dH := w.Header()
	copyHeader(resp.Header, &dH)
	dH.Add("Requested-Host", rr.Host)

	w.Write(body)
}

func defaultHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Default Value")
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func copyHeader(source http.Header, dest *http.Header){
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}