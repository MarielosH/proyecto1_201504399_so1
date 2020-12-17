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
	"time"

	"github.com/gorilla/websocket"
)

var clientes = make(map[*websocket.Conn]string)

type memoria struct {
	Total      string `json:"Total"`
	Uso        string `json:"Uso"`
	Porcentaje string `json:"Porcentaje"`
}

type datos struct {
	Ejecucion   string `json:"Ejecucion"`
	Suspendidos string `json:"Suspendidos"`
	Detenidos   string `json:"Detenidos"`
	Zombie      string `json:"Zombie"`
	Total       string `json:"Total"`
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

type objcpu struct {
	ListaProcesos []proceso `json:"ListaProcesos"`
	General       []datos   `json:"General"`
}

var actualiza = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	actualiza.CheckOrigin = func(r *http.Request) bool { return true }

	// websocket connection
	ws, err := actualiza.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}
	// returnswebsocket connection
	return ws, nil
}
func Upgrade2(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	actualiza.CheckOrigin = func(r *http.Request) bool { return true }

	// websocket connection
	ws, err := actualiza.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}
	// returnswebsocket connection
	return ws, nil
}
func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)
	ws, err := Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}
	//defer ws.Close()
	go infomemo(ws)
}

func infomemo(conn *websocket.Conn) {

	for {
		ticker := time.NewTicker(2 * time.Second)

		for t := range ticker.C {
			fmt.Printf("Cargando : %+v\n", t)
			ArchivoMemo, errores := ioutil.ReadFile("/proc/memo_201504399")
			//se puede obviar de acá
			memo := memoria{}
			errores = json.Unmarshal(ArchivoMemo, &memo)
			if errores != nil {
				fmt.Println("Error al leer el archivo memo.", errores)
				return
			}
			fmt.Println(memo)
			//hasta acá, eso era solo para leerlo
			if err := conn.WriteMessage(websocket.TextMessage, []byte(ArchivoMemo)); err != nil {
				fmt.Println(err)
				return
			}

		}

	}

	//time.Sleep(2000 * time.Millisecond)
}

func servecpu(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)
	ws, err := Upgrade2(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}
	//defer ws.Close()
	go infocpu(ws)
}
func infocpu(conn *websocket.Conn) {

	for {
		ticker := time.NewTicker(2 * time.Second)

		for t := range ticker.C {
			fmt.Printf("Cargando : %+v\n", t)
			Archivocpu, err := ioutil.ReadFile("/proc/cpu_201504399")
			//listaprocesos := []proceso{}
			objetocpu := objcpu{}
			err = json.Unmarshal(Archivocpu, &objetocpu)
			if err != nil {
				fmt.Println("Error al leer el archivo cpu.", err)
				return
			}
			fmt.Println(objetocpu.General)

			if err := conn.WriteMessage(websocket.TextMessage, []byte(Archivocpu)); err != nil {
				fmt.Println(err)
				return
			}

		}

	}
	//time.Sleep(2000 * time.Millisecond)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world")
}

func otro() {
	http.HandleFunc("/procesos", servecpu)
}
func main() {
	fmt.Println("Puerto 3000")
	fs := http.FileServer(http.Dir("./Frontend"))
	http.Handle("/", fs)
	//http.HandleFunc("/home", homePage)
	http.HandleFunc("/memo", serveWs)
	go otro()
	log.Fatal(http.ListenAndServe(":3000", nil))

}
