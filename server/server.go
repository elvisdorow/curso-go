package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CotacaoAPI struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type AwesomeAPIResultado struct {
	CotacaoAPI CotacaoAPI `json:"USDBRL"`
}

type RespostaCotacao struct {
}

func main() {
	http.HandleFunc("/cotacao", func(res http.ResponseWriter, req *http.Request) {
		apiUrl := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
		ctx, _ := context.WithTimeout(req.Context(), 200*time.Millisecond)

		clientRequest, err := http.NewRequestWithContext(ctx, "GET", apiUrl, nil)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		clientResponse, err := http.DefaultClient.Do(clientRequest)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		defer clientResponse.Body.Close()

		var cotacao AwesomeAPIResultado
		json.NewDecoder(clientResponse.Body).Decode(&cotacao)

		err = inserirCotacaoDb(&cotacao.CotacaoAPI)
		if err != nil {
			http.Error(res, "Erro ao inserir registro no banco de dados",
				http.StatusInternalServerError)
		}

		res.Header().Set("Content-type", "application/json")
		json.NewEncoder(res).Encode(&cotacao.CotacaoAPI)
	})

	http.ListenAndServe(":8080", nil)
}

func inserirCotacaoDb(cotacao *CotacaoAPI) error {
	db, err := sql.Open("sqlite3", "./banco.db")
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Millisecond)

	stmt, _ := db.Prepare(`
		insert into cotacoes (
		code, codein, name, high, low, var_bid, 
		pct_change, bid, ask, timestamp, create_date
		) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)

	_, err = stmt.ExecContext(ctx,
		cotacao.Code,
		cotacao.Codein,
		cotacao.Name,
		cotacao.High,
		cotacao.Low,
		cotacao.VarBid,
		cotacao.PctChange,
		cotacao.Bid,
		cotacao.Ask,
		cotacao.Timestamp,
		cotacao.CreateDate,
	)
	if err != nil {
		return err
	}

	return nil
}
