package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"

)

type Proceso struct {
	Id uint64
	I chan uint64
}

type Cliente struct {
	Id uint64
	I uint64
}



func butIsError(err error){
	if err != nil {
		fmt.Println(err)
		return
	}
}

func client() {
	cli, err := net.Dial("tcp", ":9999")
	butIsError(err)
	var pr Cliente
	err = gob.NewDecoder(cli).Decode(&pr)
	butIsError(err)

	proc := []Proceso{}
	channlID :=  make(chan uint64)

	procesoA := Proceso{Id: pr.Id, I: make(chan uint64)}
	proc = append(proc, procesoA)
	
	
	go func() {
		inc:=pr.I
		for {
			
			for{
				for _, proceso := range proc {
					time.Sleep(time.Second / 2)
					
					fmt.Printf("%d : %d \n",proceso.Id,inc)
					channlID<-proceso.Id
					inc++
					pr.I=inc
					
				}
			}
		}
	}()
	
	go func() {
		
		for {
		
			pr.Id= <-channlID
			
			err = gob.NewEncoder(cli).Encode(pr)
	}}()
}

func main() {
	go client()
	var input string
	fmt.Scanln(&input)
}