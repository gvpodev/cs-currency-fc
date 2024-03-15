package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", GetCurrencyInfo)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

// GetCurrencyInfo fetches the currency information from the API
func GetCurrencyInfo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", economyAPIURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}

	var info USDBRL
	if err := json.Unmarshal(body, &info); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}

	_, err = w.Write([]byte("Cotação do Dólar: " + info.USDBRL.Bid))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
}
