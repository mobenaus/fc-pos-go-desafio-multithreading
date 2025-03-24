package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BrasilAPIEndereco struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCepAPIEndereco struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type Endereco struct {
	Cep        string
	Estado     string
	Cidade     string
	Bairro     string
	Logradouro string
}

func main() {
	brasilapi := make(chan Endereco)
	viacepapi := make(chan Endereco)

	cep := "88010000"

	go callBrasilapi(brasilapi, cep)
	go callViacepapi(viacepapi, cep)

	select {
	case endereco := <-brasilapi:
		output(endereco, "BrasilAPI")
	case endereco := <-viacepapi:
		output(endereco, "ViaCep")
	case <-time.After(time.Second):
		println("timeout, nenhuma api retornou dentro de 1 segundo")
	}

}

func output(endereco Endereco, api string) {
	fmt.Printf("Api %s retornou o endereço:\n", api)
	fmt.Printf("Cep: %s\n", endereco.Cep)
	fmt.Printf("Estado: %s\n", endereco.Estado)
	fmt.Printf("Cidade: %s\n", endereco.Cidade)
	fmt.Printf("Bairro: %s\n", endereco.Bairro)
	fmt.Printf("Rua: %s\n", endereco.Logradouro)
}

func callBrasilapi(brasilapi chan Endereco, cep string) {
	req, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep)
	if err != nil {
		return
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}
	var bae BrasilAPIEndereco
	err = json.Unmarshal(res, &bae)
	if err != nil {
		return
	}

	brasilapi <- Endereco{
		Cep:        bae.Cep,
		Estado:     bae.State,
		Cidade:     bae.City,
		Bairro:     bae.Neighborhood,
		Logradouro: bae.Street,
	}

}
func callViacepapi(viacepapi chan Endereco, cep string) {
	req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}
	var vae ViaCepAPIEndereco
	err = json.Unmarshal(res, &vae)
	if err != nil {
		return
	}

	viacepapi <- Endereco{
		Cep:        vae.Cep,
		Estado:     vae.Estado,
		Cidade:     vae.Localidade,
		Bairro:     vae.Bairro,
		Logradouro: vae.Logradouro,
	}

}
