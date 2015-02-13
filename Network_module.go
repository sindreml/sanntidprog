//Network module
package main

import ("fmt"
		"net"
		"time"
		"encoding/json")
		
		

func main(){

	NetworkModule()
}


type StatePackage struct {
	Name string
	Level int
	AttemptNumber int
	Hjemmebane string
	Bortebane string
}

// Innmaten i NetworkModeulen er laget for å teste selve modulen. Når endelig modul implementeres, husk å bruke select hver gang det brukes en kanal 
//ellers vil programmet bli stående å vente hvis man ikke kan lese / skrive til kanalen
//Obs! Når man printer en struct vil ikke nøklene vises, kun verdiene som nøklene peker til. Dette er sannsynligvis ikke noe problem.

func NetworkModule(){ 
	sendObject:= StatePackage{"Sindre",112,3432,"brille","balle"}
	sendChan := make(chan StatePackage,1)
	recieveChan:= make(chan StatePackage,1)
	go sendToServer("","20021",sendChan)
	go readFromServer("","20021",recieveChan)
	sendChan <- sendObject	
	
	for {
		select {
		case testObject:= <- recieveChan:
			fmt.Println("det var mulig å lese fra recieveChan i main. Vi leser:",testObject)		
		case sendChan <- sendObject:	
			
		default:
			time.Sleep(50*time.Millisecond)		
		}

		time.Sleep(2000*time.Millisecond)
	}
}


func readFromServer(ipAddress string, portNumber string, recieveChan chan StatePackage) { 		
	bufferToRead := make([] byte, 1024)
	UDPadr, err:= net.ResolveUDPAddr("udp",ipAddress+":"+portNumber)

	if err != nil {
                fmt.Println("error resolving UDP address on ", portNumber)
                fmt.Println(err)
                return
        }
    
    readerSocket ,err := net.ListenUDP("udp",UDPadr)
    
    if err != nil {
            fmt.Println("error listening on UDP port ", portNumber)
            fmt.Println(err)
            return
	}
	
	for {
		n,UDPadr, err := readerSocket.ReadFromUDP(bufferToRead)
        fmt.Println("inne i readFromServer")

       	if err != nil {
            fmt.Println("error reading data from connection")
            fmt.Println(err)
            return
        }
        
        fmt.Println("got message from ", UDPadr, " with n = ", n)

       	if n > 0 {
           	fmt.Println("printer melding vi leste over UDP",json2struct(bufferToRead[0:n]))  
            structObject := json2struct(bufferToRead[0:n])
           	recieveChan <- structObject
        }
	}
}

func sendToServer(ipAddress string, portNumber string, sendChan chan StatePackage){
	socketSend, err1 := net.Dial("udp", ipAddress+":"+portNumber)
	if err1 != nil {
                fmt.Println("error listening on UDP port ", portNumber)
                fmt.Println(err1)
                return
	}
	
	for {
		time.Sleep(5000*time.Millisecond)
		select{
		case <- sendChan:
			packageToSend :=  <- sendChan
			jsonFile := struct2json(packageToSend)
			socketSend.Write(jsonFile)
			fmt.Println("jsonFile er sendt over UDP")
		default:
			fmt.Println("sendChan var tom")
			time.Sleep(50*time.Millisecond)
		}			
	}
} 

func struct2json(packageToSend StatePackage) [] byte {
	jsonObject, _ := json.Marshal(packageToSend)
	return jsonObject
}

func json2struct(jsonObject []byte) StatePackage{
	structObject := StatePackage{}
	json.Unmarshal(jsonObject, &structObject)
	return structObject
}
