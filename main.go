package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var lock sync.Mutex

var client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func main() {
	url := flag.String("url", "", "endereço a ser testado")
	total := flag.Int("requests", 0, "qtde total de requisições")
	concurrent := flag.Int("concurrency", 0, "qtde de requisições concorrentes")

	flag.Parse()

	if *url == "" {
		panic("-url não pode ser vazia")
	}
	if *total <= 0 || *concurrent <= 0 {
		panic("-requests e -concurrency devem ser maiores que 0")

	}

	if *total < *concurrent {
		panic("-requests deve ser maior ou igual a -concurrency")
	}

	var failedRequests = 0
	var rounds = *total / *concurrent
	var lastRound = *total % *concurrent
	var mapOfResponses = make(map[int]int)

	wg := sync.WaitGroup{}
	wg.Add(*total)

	start := time.Now()

	for r := 0; r < rounds; r++ {
		for i := 0; i < *concurrent; i++ {
			go doRequest(*url, &failedRequests, &mapOfResponses, &wg)
		}
	}

	if lastRound > 0 {
		for i := 0; i < lastRound; i++ {
			go doRequest(*url, &failedRequests, &mapOfResponses, &wg)
		}
	}

	wg.Wait()

	executionTime := time.Since(start)

	fmt.Printf("%d requisições executadas com sucesso de um total de %d tentativas.\n", *total-failedRequests, *total)
	fmt.Printf("Tempo de execução: %s\n", executionTime)

	if failedRequests == *total {
		fmt.Println("Todas as requisições falharam.")
	} else {
		fmt.Println("\n\nRelatório detalhado de status code por número de requisições:")
		for key, value := range mapOfResponses {
			fmt.Printf("Código HTTP: %d, Quantidade: %d\n", key, value)
		}
	}
}

func doRequest(url string, failedRequests *int, responses *map[int]int, wg *sync.WaitGroup) {
	lock.Lock()
	defer lock.Unlock()
	println("Iniciou uma requisição...")
	start := time.Now()

	req, err := client.Get(url)

	defer wg.Done()

	if err != nil {
		fmt.Printf("Requisição falhou com erro %v\n", err)
		*failedRequests++
		return
	} else {
		req.Close = true
	}

	code := req.StatusCode

	println("Requisição finalizada com status code ", code)

	value, ok := (*responses)[code]
	if ok {
		(*responses)[code] = value + 1
	} else {
		(*responses)[code] = 1
	}

	fmt.Printf("Requisição finalizada em %s\n", time.Since(start))

	defer req.Body.Close()
}
