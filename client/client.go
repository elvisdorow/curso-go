package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type CotacaoAtual struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 300*time.Millisecond)

	apiUrl := "http://localhost:8080/cotacao"
	clientRequest, err := http.NewRequestWithContext(ctx, "GET", apiUrl, nil)
	tratarErro(err)

	clientResponse, err := http.DefaultClient.Do(clientRequest)
	tratarErro(err)

	defer clientResponse.Body.Close()

	var cotacao CotacaoAtual
	json.NewDecoder(clientResponse.Body).Decode(&cotacao)

	arquivo, err := os.Create("cotacao.txt")
	tratarErro(err)

	defer arquivo.Close()

	arquivo.WriteString(fmt.Sprintf("Dólar: %v", cotacao.Bid))
}

func tratarErro(err error) {
	if err != nil {
		panic(err)
	}
}
