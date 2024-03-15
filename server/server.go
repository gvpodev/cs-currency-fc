package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type CustomTime struct {
	time.Time
}

const ctLayout = "2006-01-02 15:04:05"

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
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	req, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	var info USDBRL
	if err := json.Unmarshal(res, &info); err != nil {
		log.Fatal(err)
	}

	log.Printf("Cotação do Dólar: %v", info.USDBRL.Bid)
}
