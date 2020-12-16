package main

import (
	"C"
	"fmt"
	"net/http"
)

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gorilla/websocket"
)

var clientes = make(map[*websocket.Conn]string)

var actualiza = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	ws, err := actualiza.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	defer ws.Close()

	reader(ws)
}

func reader(conn *websocket.Conn) {

	for {

		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("error: %v", err)
			delete(clientes, conn)
			break
		}

		fmt.Println(string(p))
		clientes[conn] = string(p)
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

type memoria struct {
	Total      string `json:"Total"`
	Uso        string `json:"Uso"`
	Porcentaje string `json:"Porcentaje"`
}

type hijo struct {
	PID     string `json:"PID"`
	PROCESO string `json:"PROCESO"`
	ESTADO  string `json:"ESTADO"`
	MEMORIA string `json:"MEMORIA"`
	USUARIO string `json:"USUARIO"`
}

type proceso struct {
	PID     string `json:"PID"`
	PROCESO string `json:"PROCESO"`
	ESTADO  string `json:"ESTADO"`
	MEMORIA string `json:"MEMORIA"`
	USUARIO string `json:"USUARIO"`
	HIJOS   []hijo `json:"HIJOS"`
}

func infomemo() {
	ArchivoMemo, errores := ioutil.ReadFile("/proc/memo_201504399")
	memo := memoria{}
	errores = json.Unmarshal(ArchivoMemo, &memo)
	if errores != nil {
		fmt.Println("Error al leer el archivo memo.", errores)
		return
	}
	fmt.Println(memo.Total)
	fmt.Println(memo.Uso)
	fmt.Println(memo.Porcentaje)

}

func infocpu() {
	Archivocpu, err := ioutil.ReadFile("/proc/cpu_201504399")
	listaprocesos := []proceso{}
	err = json.Unmarshal(Archivocpu, &listaprocesos)
	if err != nil {
		fmt.Println("Error al leer el archivo cpu.", err)
		return
	}
	//fmt.Println(listaprocesos)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
	fmt.Println("Hola mundo")
	infomemo()
	infocpu()
	setupRoutes()
}
