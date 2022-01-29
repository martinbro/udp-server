package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// remoteAddress ligger som global variabel
var remoteAddress *net.UDPAddr

type myConn struct {
	ch     chan []byte
	rate   chan int64
	msg    chan []byte
	dublex bool
	// remoteAddress *net.UDPAddr //adresse til skib
	// serverconn    *net.UDPAddr //forbindelse mellem server og skib
}

func (c myConn) handler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil) //returnerer en upgradet sucket
	if err != nil {
		log.Println(err)
		return
	}

	//spinder er gorutine op, der skriver til en websocket (app´en)
	go func(ws *websocket.Conn) {

		for {
			q := <-c.ch
			if err := ws.WriteMessage(1, q); err != nil {
				fmt.Fprintln(os.Stderr, "Fejl i writer: ", err)
				return
			}
		}
	}(ws)

	if c.dublex {
		//spinder er gorutine op, der læser fra en websocket (app´en)
		go func(ws *websocket.Conn, c myConn) {
			for {
				messageType, p, err := ws.ReadMessage()
				if err != nil {
					log.Println("Fejl i dasboard reader: ", err)
					return
				}

				if messageType == websocket.TextMessage {
					label := string(p[0:5])
					s := string(p[5:])
					fmt.Printf("TextMessage: %s %s %s\n", p, s, label)
					switch label {
					case "esp;b":
						if data, err := strconv.Atoi(s); err == nil {
							d := int64(data)
							c.rate <- d
						}
					default:
						c.msg <- p
					}
				} else {
					fmt.Printf("Anden medd.....: %s\n", p)
				}
			}
		}(ws, c)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	//initialiser chan
	bno := myConn{ch: make(chan []byte), msg: make(chan []byte), rate: make(chan int64)}
	gps := myConn{ch: make(chan []byte), msg: make(chan []byte)}
	ws2 := myConn{ch: make(chan []byte), msg: make(chan []byte), rate: bno.rate, dublex: true}

	fmt.Println("bno", bno.rate)
	fmt.Println("gps", gps.rate)

	go setupUDP(bno, "192.168.137.1:8081")
	go setupUDP(gps, "192.168.137.1:8082")
	go setupUDP(ws2, "192.168.137.1:8083")

	http.HandleFunc("/ws1bno", bno.handler)
	http.HandleFunc("/ws1gps", gps.handler)
	http.HandleFunc("/ws2", ws2.handler)

	fmt.Println("Følg linket (Alt + click):", "http://192.168.137.1:8000")
	// //Midlertidig testing
	// go varierRate(bno.rate)
	bno.rate <- 20 //Default

	Openbrowser("http://192.168.137.1:8000")
	http.ListenAndServe("192.168.137.1:8000", nil) //Blocking call
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//Etablerer "BNO"-data forbindelsen fra server til Dashboard

// func sendBesked(ch chan string) {
// 	for {
// 		time.Sleep(2 * time.Second)
// 		ch <- "besked"
// 	}
// }

// func varierRate(rate chan int64) {
// 	for {
// 		speed := rand.Int63n(900) + 1
// 		rate <- speed
// 		time.Sleep(2 * time.Second)
// 	}
// }

func newConn(s string) *net.UDPConn {
	//Forbinder til UDP-socket
	udpAddr, _ := net.ResolveUDPAddr("udp", s)
	serverConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Print("\nKunne ikke forbinde til mobilt hotspot!\nHar du husket at tænde det mobile hotspot?\n")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
	}
	remoteAddress = udpAddr
	return serverConn
}

func setupUDP(con myConn, url string) {

	serverConn := newConn(url)
	defer serverConn.Close()

	//Initialiserer
	var rateVal int64 = 10
	t0 := time.Now()
	buffer := make([]byte, 50) //NB 50 tegn i buffer
	i := 0

	for { //starter i en uendelig løkke
		_, remoteAdd, err := serverConn.ReadFromUDP(buffer)
		// n, remoteAdd, err := serverConn.ReadFromUDP(buffer)
		if err != nil {
			println("FEJL i ", err)
		}
		if remoteAddress != remoteAdd {
			remoteAddress = remoteAdd
		}
		select {
		case r := <-con.rate:
			rateVal = r
		case msg := <-con.msg:
			fmt.Println("Sender beskeden!")
			serverConn.WriteToUDP(msg, remoteAddress)
		default:
		}
		if time.Since(t0) > time.Duration(rateVal*int64(time.Millisecond)) {
			t0 = time.Now()
			// navn := string(buffer[0:3])

			// fmt.Printf("\nClient: %s - %s - %d", navn, buffer, n)
			// if navn == "pin" {
			// 	fmt.Printf("\nClient: %s - %s - %d - %s", navn, buffer, n, remoteAddress)
			// }
			select {
			case con.ch <- buffer:

			default:
			}
		}
		i++
		if i%10 == 0 {
			print("+")
		} else {
			print(".")

		}
	}
}

// //EXPRIMENTELT
// type myConn1 struct {
// 	ch            chan []byte
// 	rate          chan int64
// 	msg           chan []byte
// 	dublex        bool
// 	ws            *websocket.Conn
// 	err           error
// 	remoteAddress *net.UDPAddr //adresse til skib
// 	serverconn    *net.UDPAddr //forbindelse mellem server og skib

// }

// //Metode
// func (c *myConn1) getWebsocket(w http.ResponseWriter, r *http.Request) {
// 	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
// 	c.ws, c.err = upgrader.Upgrade(w, r, nil) //returnerer en upgradet sucket
// }
