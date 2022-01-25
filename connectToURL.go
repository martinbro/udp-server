package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func Openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

// func N(ch chan string, rate chan int64) {
// 	//Forbinder til UDP-socket
// 	udpAddr, _ := net.ResolveUDPAddr("udp", "192.168.137.1:8081")
// 	serverConn, err := net.ListenUDP("udp", udpAddr)
// 	if err != nil {
// 		fmt.Print("\nKunne ikke forbinde til mobilt hotspot.\n\nHar du husket at tænde for det mobile hotspot?\n")
// 		scanner := bufio.NewScanner(os.Stdin)
// 		scanner.Scan()
// 		return
// 	}
// 	defer serverConn.Close()

// 	//Initialiserer
// 	var rateVal int64
// 	rateVal = 500
// 	t0 := time.Now()
// 	buffer := make([]byte, 100)
// 	last := 0
// 	i := 0

// 	for {
// 		n, remoteAddress, err := serverConn.ReadFromUDP(buffer)
// 		if err != nil {
// 			println("FEJL", err)
// 		}
// 		select {
// 		case r := <-rate:
// 			rateVal = r
// 		default:
// 		}
// 		select {
// 		case msg := <-ch:
// 			serverConn.WriteToUDP([]byte(msg), remoteAddress)
// 		default:
// 		}
// 		// fmt.Println("remoteAddress", remoteAddress.IP, remoteAddress.Port, remoteAddress.Zone)
// 		if time.Since(t0) > time.Duration(rateVal*int64(time.Millisecond)) {
// 			t0 = time.Now()

// 			nr := (buffer[4 : n-8])
// 			x, err := strconv.Atoi(string(nr))
// 			if err == nil {

// 				fmt.Print("\nClient ", remoteAddress, ":", x, x-last)
// 				fmt.Printf(" %s", nr)
// 				last = x
// 				i = 0
// 			}
// 		}
// 		i++
// 		if i%10 == 0 {
// 			print("+")
// 		} else {
// 			print(".")

// 		}
// 	}
// }
