package Network

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"
)

func (network *Network)SendPacket(message []byte, ip string, port int) {
	port_s := strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", ip+":"+port_s, 3 * time.Second)
	if err != nil {
		if ip == "121.15.171.86" || ip == "127.0.0.1" {
			return
		}
		fmt.Println("dial failed, err", err)
		go network.updateNodeStatus(ip, port , 0)
		return
	}
	defer conn.Close()
	go network.updateNodeStatus(ip, port , 1)
	data, err := Encode(message)
	if err != nil {
		fmt.Println("encode msg failed, err:", err)
		return
	}
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
}

func (network *Network) SendPacketWithSavedConn(message []byte, ip string, port int, receiver NodeID) {
	var conn net.Conn
	conn_m, ok:= network.NodeConnList[receiver]
	if ok && conn_m != nil {
		conn = conn_m
	} else {
		var err error
		port_s := strconv.Itoa(port)
		conn, err = net.DialTimeout("tcp", ip+":"+port_s, 3 * time.Second)
		if err != nil {
			//fmt.Println("dial failed, err", err)
			//go network.updateNodeStatus(ip, port , 0)
			return
		}
		network.NodeConnList[receiver] = conn
	}

	//go network.updateNodeStatus(ip, port , 1)
	data, err := Encode(message)
	if err != nil {
		fmt.Println("encode msg failed, err:", err)
		return
	}
	_, err = conn.Write(data)
	if err != nil {
		//fmt.Println("send msg failed, err:", err)
		delete(network.NodeConnList, receiver)
		network.SendPacketWithSavedConn(message, ip, port, receiver)
		return
	}
}

func (network *Network)SendPacketAndGetAns(message []byte, ip string, port int) []byte{
	port_s := strconv.Itoa(port)
	conn, err := net.DialTimeout("tcp", ip+":"+port_s, 3 * time.Second)
	if err != nil {
		fmt.Println("dial failed, err", err)
		go network.updateNodeStatus(ip, port , 0)
		return nil
	}
	defer conn.Close()
	go network.updateNodeStatus(ip, port , 1)
	data, err := Encode(message)
	if err != nil {
		fmt.Println("encode msg failed, err:", err)
		return nil
	}
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return nil
	}

	msg, err := Decode(conn)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		fmt.Println("decode msg failed, err:", err)
		return nil
	}
	return msg
}

func (network *Network) SendToAll(message []byte) {
	network.mutex.RLock()
	tmpList := network.NodeList
	for k, x := range tmpList {
		network.SendPacketWithSavedConn(message, x.IP, x.PORT, k)
	}
	network.mutex.RUnlock()
}

func (network *Network) SendToNeighbor(message []byte) {
	network.mutex.RLock()
	tmpList := network.NodeList
	temp := network.MyNodeInfo.ID
	for k, x := range tmpList {
		if x.ID == temp {
			continue
		}
		network.SendPacketWithSavedConn(message, x.IP, x.PORT, k)
	}
	network.mutex.RUnlock()
}

func (network *Network) SendToOne(message []byte, receiver NodeID) {
	network.mutex.RLock()
	ip := network.NodeList[receiver].IP
	port := network.NodeList[receiver].PORT
	network.SendPacketWithSavedConn(message, ip, port, receiver)
	network.mutex.RUnlock()
}

func (network *Network) SendMessage(message []byte, receiver NodeID) {
	if receiver == 0 {
		network.SendToAll(message)
	} else if receiver == 1 {
		network.SendToNeighbor(message)
	} else {
		network.SendToOne(message, receiver)
	}
}

func (network *Network) updateNodeStatus( ip string , port int, status int){
	network.mutex.Lock()
	nodeId := GetNodeId(ip, port)
	node, ok :=network.NodeList[nodeId]
	if !ok {
		return
	}
	node.Status = status
	network.NodeList[nodeId] = node
	network.mutex.Unlock()
}