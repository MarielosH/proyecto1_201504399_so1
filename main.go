package main

import (
	"C"
	"fmt"
	"net/http"
)

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

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

type cpuStruct struct {
	Porcentaje float64 `json:"porcentaje"`
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
func Upgrade3(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
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

			if errores != nil {
				fmt.Println("Error al leer el archivo memo.", errores)
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, []byte(ArchivoMemo)); err != nil {
				fmt.Println(err)
				return
			}

		}

	}

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
func servecpu2(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)
	ws, err := Upgrade3(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}
	//defer ws.Close()
	go getCpuInfo(ws)
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

func getCpuInfo(conn *websocket.Conn) {
	var prevIdleTime, prevTotalTime uint64
	var cpuUsage = 0.0

	for {
		ticker := time.NewTicker(2 * time.Second)

		for t := range ticker.C {
			fmt.Printf("Cargando : %+v\n", t)
			file, err := os.Open("/proc/stat")
			if err != nil {
				log.Fatal(err)
			}
			scanner := bufio.NewScanner(file)
			scanner.Scan()
			firstLine := scanner.Text()[5:] // get rid of cpu plus 2 spaces
			file.Close()
			if err := scanner.Err(); err != nil {
				log.Fatal(err)
			}
			split := strings.Fields(firstLine)
			idleTime, _ := strconv.ParseUint(split[3], 10, 64)
			totalTime := uint64(0)
			for _, s := range split {
				u, _ := strconv.ParseUint(s, 10, 64)
				totalTime += u
			}

			deltaIdleTime := idleTime - prevIdleTime
			deltaTotalTime := totalTime - prevTotalTime
			cpuUsage = (1.0 - float64(deltaIdleTime)/float64(deltaTotalTime)) * 100.0
			fmt.Printf(" %6.3f\n", cpuUsage)

			prevIdleTime = idleTime
			prevTotalTime = totalTime

			cpuObj := &cpuStruct{math.Round(cpuUsage*100) / 100}
			jsonResponse, errorjson := json.Marshal(cpuObj)
			if errorjson != nil {
				fmt.Println("Error al leer el archivo cpu.", errorjson)
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, []byte(jsonResponse)); err != nil {
				fmt.Println(err)
				return
			}
		}
	}

}

func otro() {
	http.HandleFunc("/procesos", servecpu)
}
func otro2() {
	http.HandleFunc("/cpuinfo", servecpu2)
}
func main() {
	fmt.Println("Puerto 3000")
	fs := http.FileServer(http.Dir("./Frontend"))
	http.Handle("/", fs)
	http.HandleFunc("/memo", serveWs)
	go otro()
	go otro2()
	log.Fatal(http.ListenAndServe(":3000", nil))

}
