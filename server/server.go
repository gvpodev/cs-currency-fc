package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
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
	db, err := sql.Open("sqlite3", "currency.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = executeMigrations(db, "./migrations")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", GetCurrencyInfo)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertCurrencyInfo(info *CurrencyInfo) error {
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

	_, err = stmt.Exec(info.Code, info.Codein, info.Name, info.High, info.Low, info.VarBid, info.PctChange, info.Bid,
		info.Ask, info.Timestamp, info.CreateDate.Time)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// GetCurrencyInfo fetches the currency information from the API
func GetCurrencyInfo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
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

	err = InsertCurrencyInfo(&info.USDBRL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func executeMigrations(db *sql.DB, dirPath string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			fmt.Printf("Aplicando migração: %s\n", file.Name())

			content, err := ioutil.ReadFile(filepath.Join(dirPath, file.Name()))
			if err != nil {
				return err
			}

			tx, err := db.Begin()
			if err != nil {
				return err
			}

			if _, err := tx.Exec(string(content)); err != nil {
				tx.Rollback()
				return err
			}

			err = tx.Commit()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
