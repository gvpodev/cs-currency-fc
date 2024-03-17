package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"server/migrations"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CustomTime struct {
	time.Time
}

const (
	ctLayout      = "2006-01-02 15:04:05"
	economyAPIURL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
)

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b[1 : len(b)-1])
	t, err := time.Parse(ctLayout, s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

type USDBRL struct {
	USDBRL CurrencyInfo `json:"USDBRL"`
}

type CurrencyInfo struct {
	Code       string     `json:"code"`
	Codein     string     `json:"codein"`
	Name       string     `json:"name"`
	High       string     `json:"high"`
	Low        string     `json:"low"`
	VarBid     string     `json:"varBid"`
	PctChange  string     `json:"pctChange"`
	Bid        string     `json:"bid"`
	Ask        string     `json:"ask"`
	Timestamp  string     `json:"timestamp"`
	CreateDate CustomTime `json:"create_date"`
}

func main() {
	migrateDB := flag.Bool("migratedb", false, "Set true to execute database migration")
	flag.Parse()
	if *migrateDB {
		migrations.Execute()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", GetInfoHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

// GetInfoHandler fetches the currency information from the API
func GetInfoHandler(w http.ResponseWriter, _ *http.Request) {
	info, err := getCurrencyInfo()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Erro ao buscar informações de cotação: %v\n", err)
		_, _ = w.Write([]byte("Erro ao buscar informações de cotação"))

		return
	}

	err = json.NewEncoder(w).Encode(&info.USDBRL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Erro ao processar resposta: %v\n", err)
		_, _ = w.Write([]byte("Erro ao buscar informações de cotação"))

		return
	}

	err = insertCurrencyInfo(&info.USDBRL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Printf("Erro ao inserir resposta: %v\n", err)
		return
	}
}

func getCurrencyInfo() (*USDBRL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", economyAPIURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var info *USDBRL
	if err = json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	return info, err
}

func insertCurrencyInfo(info *CurrencyInfo) error {
	db, err := sql.Open("sqlite3", "currency.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO currency (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date)" +
		" VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = stmt.ExecContext(ctx, info.Code, info.Codein, info.Name, info.High, info.Low, info.VarBid, info.PctChange, info.Bid,
		info.Ask, info.Timestamp, info.CreateDate.Time)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
