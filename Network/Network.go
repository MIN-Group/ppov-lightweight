package Network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"

	"ppov/ConfigHelper"
)

type NodeID = uint64
type NodeInfo struct {
	IP   string
	PORT int
	ID   NodeID
	Status int 
}

type Network struct {
	MyNodeInfo NodeInfo
	NodeList   map[NodeID]NodeInfo
	ServicePort	int
	CB         func([]byte, net.Conn)

	//TCP长连接
	ConnSliceBC []net.Conn
	ConnSliceMIN []net.Conn
	NodeConnList map[NodeID]net.Conn

	//锁
	mutex sync.RWMutex
}

func IPToValue(strIP string) uint32 {
	var a [4]uint32
	temp := strings.Split(strIP, ".")
	for i, x := range temp {
		t, err := strconv.Atoi(x)
		if err != nil {
			fmt.Println(err)
		}
		a[i] = uint32(t)
	}
	var ret uint32 = (a[0] << 24) + (a[1] << 16) + (a[2] << 8) + a[3]
	return ret
}

func GetNodeId(ip string, port int) NodeID {
	var id uint32 = IPToValue(ip)
	var nid NodeID = NodeID(id) << 32
	nid += NodeID(port)
	return nid
}

func (network *Network) SetConfig(config ConfigHelper.Config) {
	network.MyNodeInfo.IP = config.MyAddress.IP
	network.MyNodeInfo.PORT = config.MyAddress.Port
	network.MyNodeInfo.ID = GetNodeId(network.MyNodeInfo.IP, network.MyNodeInfo.PORT)
	network.ServicePort = config.ServicePort
	network.NodeList = make(map[NodeID]NodeInfo)
	for _, x := range config.WorkerList {
		temp := GetNodeId(x.IP, x.Port)
		_, ok := network.NodeList[temp]
		if !ok {
			var nodelist NodeInfo
			nodelist.IP = x.IP
			nodelist.PORT = x.Port
			nodelist.ID = temp
			network.NodeList[temp] = nodelist
		}
	}
	for _, x := range config.WorkerCandidateList {
		temp := GetNodeId(x.IP, x.Port)
		_, ok := network.NodeList[temp]
		if !ok {
			var nodelist NodeInfo
			nodelist.IP = x.IP
			nodelist.PORT = x.Port
			nodelist.ID = temp
			network.NodeList[temp] = nodelist
		}
	}
	for _, x := range config.VoterList {
		temp := GetNodeId(x.IP, x.Port)
		_, ok := network.NodeList[temp]
		if !ok {
			var nodelist NodeInfo
			nodelist.IP = x.IP
			nodelist.PORT = x.Port
			nodelist.ID = temp
			network.NodeList[temp] = nodelist
		}
	}
}

func (network *Network) SetCB(cb func([]byte, net.Conn)) {
	network.CB = cb
}

// Encode 将消息编码
func Encode(message []byte) ([]byte, error) {
	// 读取消息的长度，转换成int32类型（占4个字节）
	var length = int32(len(message))
	var pkg = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, message)
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// Decode 解码消息
func Decode(conn net.Conn) ([]byte, error) {
	// 读取消息的长度
	buf := make([]byte, 4)
	_, err := io.ReadFull(conn, buf) // 读取前4个字节的数据
	if err != nil {
		return make([]byte, 0), err
	}
	length := binary.LittleEndian.Uint32(buf)
	// 读取真正的消息数据
	pack := make([]byte, length)
	_, err = io.ReadFull(conn, pack)
	if err != nil {
		return make([]byte, 0), err
	}
	return pack, nil
}
