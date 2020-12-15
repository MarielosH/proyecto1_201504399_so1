package main

import (
	"C"
	"fmt"
)

type Proceso struct {
	PID     string    `json:"PID"`
	PROCESO string    `json:"PROCESO"`
	ESTADO  string    `json:"ESTADO"`
	HIJOS   []Proceso `json:"HIJOS"`
}

type memoria struct {
	Total      string `json:"Total"`
	Uso        string `json:"Uso"`
	Porcentaje string `json:"Porcentaje"`
}

func main() {
	fmt.Println("Hola mundo")
}
