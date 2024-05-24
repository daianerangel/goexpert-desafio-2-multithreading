package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Address struct {
	Cep          string `json:"cep"`
	Logradouro   string `json:"logradouro,omitempty"`
	Bairro       string `json:"bairro,omitempty"`
	Localidade   string `json:"localidade,omitempty"`
	Uf           string `json:"uf,omitempty"`
	Street       string `json:"street,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
	State        string `json:"state,omitempty"`
	City         string `json:"city,omitempty"`
	Source       string `json:"-"`
}

func fetchFromBrasilAPI(ctx context.Context, cep string, ch chan<- Address) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var address Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err == nil {
		address.Source = "BrasilAPI"
		ch <- address
	}
}

func fetchFromViaCep(ctx context.Context, cep string, ch chan<- Address) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var address Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err == nil {
		address.Source = "ViaCep"
		ch <- address
	}
}

func main() {
	cep := "23093010"
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	ch := make(chan Address, 2) // Channel to receive addresses

	go fetchFromBrasilAPI(ctx, cep, ch)
	go fetchFromViaCep(ctx, cep, ch)

	select {
	case address := <-ch:
		json, _ := json.Marshal(address)
		fmt.Printf("Fastest response from %s: %+v\n", address.Source, string(json))
	case <-ctx.Done():
		fmt.Println("Timeout error: No response within 1 second.")
	}
}
