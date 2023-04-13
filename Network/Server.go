package Network

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

func (network *Network) deleteConn(conn net.Conn)error{
	if conn==nil{
		fmt.Println("conn is nil")
		return errors.New("conn is nil")
	}
	for i:= 0;i<len(network.ConnSliceBC);i++{
		if network.ConnSliceBC[i]==conn {
			network.ConnSliceBC = append(network.ConnSliceBC[:i],network.ConnSliceBC[i+1:]...)
			break
		}
	}
	return nil
}

// for blockchain
func (network *Network) HandleConnection(conn net.Conn) {
	defer func(){
		conn.Close()
		network.deleteConn(conn)
	}()
	for {
		msg, err := Decode(conn)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode msg failed, err:", err)
			return
		}
		network.CB(msg, conn)
	}
}

func (network *Network) Start() {
	port := strconv.Itoa(network.MyNodeInfo.PORT)
	tcpAdd,err:= net.ResolveTCPAddr("tcp", "0.0.0.0"+":"+port)
	if err!=nil{
		log.Fatal(err)
		return
	}
	tcpListener,err:=net.ListenTCP("tcp4",tcpAdd)
	if err!=nil{
		log.Fatal(err)
		return
	}
	defer tcpListener.Close()

	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		network.ConnSliceBC = append(network.ConnSliceBC, conn)
		go network.HandleConnection(conn)
	}
}
