package main

import (
	"encoding/gob"
	"fmt"
	"sort"
	"net"
	"time"
)

type Proceso struct {
	Id uint64
	I  chan uint64
}

type Servidor  struct {
	Id uint64
	I  uint64
}

func (proc *Proceso) incremento_numero(numero uint64) {
	i := uint64(numero)
	for {
		proc.I <- i
		i++
	}
}

func butIsError(err error) {
	if err != nil {
		panic(err)
	}
}

func runServidor() {
	server, err := net.Listen("tcp", ":9999")
	procesos := []Proceso{}
	const maximoProcesos = 5
	var j uint64
	for j = 0; j < maximoProcesos; j++ {
		proc := Proceso{Id: j, I: make(chan uint64)}
		procesos = append(procesos, proc)
		go proc.incremento_numero(0)
	}
	/*if err != nil {
		fmt.Println(err)
		return
	}*/

	butIsError(err)

	go func() {
		for {
			time.Sleep(time.Second / 2)
			sort.SliceStable(procesos, func(i, j int) bool {
				return procesos[i].Id < procesos[j].Id
			})
			for _, proceso := range procesos {
				chanel := <-proceso.I
				fmt.Printf("%d: %d \n", proceso.Id, chanel)
			}
			fmt.Println("----------")
		}
	}()

	for {
		cli, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleClient(cli, &procesos)
	}
}

func handleClient(c net.Conn, proc *[]Proceso) {
	var proceso Proceso
	//var response string 
	defer c.Close() 
	if len(*proc) > 0 {
		proceso = (*proc)[0]
		i := <-proceso.I
		inf := Servidor{Id: proceso.Id, I: i}
		err := gob.NewEncoder(c).Encode(inf)
		butIsError(err)
		remove(proc, 0)
	}

	
	
	var ser Servidor
	
	for {
		err := gob.NewDecoder(c).Decode(&ser)
		if err != nil {
	
			psend := Proceso{Id:ser.Id, I: make(chan uint64)}
			go psend.incremento_numero(ser.I)
			add(proc, psend)
			return
		}

	}
}

func remove(slice *[]Proceso, s int) {
	ss := *slice
	ss = append(ss[:s], ss[s+1:]...)
	*slice = ss
}

func add(slice *[]Proceso, s Proceso) {
	ss := *slice
	ss = append(ss, s)
	*slice = ss
}

func main() {
	go runServidor()

	var input string
	fmt.Scanln(&input)
}
