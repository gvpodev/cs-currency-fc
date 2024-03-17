package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const serverURL = "http://localhost:8080/cotacao"

type ClientResponse struct {
	Bid string `json:"bid"`
}

func main() {
	fmt.Println("Iniciando cliente...")

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", serverURL, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var cResponse *ClientResponse
	err = json.Unmarshal(body, &cResponse)
	if err != nil {
		log.Printf("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("response: %q", body)
		return
	}

	file, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	insert := fmt.Sprintf("Dólar %s \n", cResponse.Bid)
	_, err = file.Write([]byte(insert))
	if err != nil {
		panic(err)
	}

	fmt.Println("Cotação inserida com sucesso")
}
