package Client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
	"ppov/KeyManager"
)

func HTTPClient(){
	fmt.Println("请选择需要的函数: \n 1.获取当前区块高度信息\n2.创建账户\n3.创建账户测试")
	var choice int
	fmt.Scanln(&choice)
	switch choice {
	case 1:
		url := "http://10.0.0.19/GetCurrentHeight"
		contentType := "application/json;charset=utf-8"
		body := bytes.NewBuffer(nil)

		resp, err := http.Post(url, contentType, body)
		if err != nil {
			log.Println("Post failed:", err)
			return
		}

		defer resp.Body.Close()

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Read failed:", err)
			return
		}

		fmt.Println(string(content))
	case 2:
		url := "http://10.0.0.19:8010/CreateAccount"
		contentType := "application/json;charset=utf-8"
		var s CreatAccountMsg
		now := time.Now()
		s.Timestamp = fmt.Sprintf("%d", now.Unix())

		var km KeyManager.KeyManager
		km.Init()
		km.GenRandomKeyPair()
		s.Address = km.GetAddress()
		s.Pubkey = km.GetPubkey()
		s.Sig = ""
		temp2, _ := json.Marshal(s)
		s.Sig, _= km.Sign(KeyManager.GetHash(temp2))
		pri, _ := km.GetPriKey()
		fmt.Println("私钥:", pri)
		fmt.Println("公钥", s.Pubkey)
		fmt.Println("地址", s.Address)
		KeyTxtWrite("PubKey", s.Address, s.Pubkey) //写入公钥
		KeyTxtWrite("PriKey", s.Address, pri)      //写入私钥
		data, err := json.Marshal(s)
		if err != nil {
			fmt.Println(err)
		}

		body := bytes.NewBuffer(data)
		resp, err := http.Post(url, contentType, body)
		if err != nil {
			log.Println("Post failed:", err)
			return
		}

		defer resp.Body.Close()

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Read failed:", err)
			return
		}

		fmt.Println(string(content))
	case 3:
		var cs int
		var gs int
		fmt.Println("请输入client数量:")
		fmt.Scan(&cs)
		fmt.Println("请输入每个client发送数量:")
		fmt.Scan(&gs)

		var num int32 = 0
		var wg sync.WaitGroup
		all := cs * 2
		wg.Add(all)

		var messages []CreatAccountMsg
		mutex := sync.Mutex{}
		for i := 0; i < cs*2; i++ {

			go func() {
				defer wg.Done()
				for j := 0; j < gs; j++ {
					var s CreatAccountMsg
					now := time.Now()
					s.Timestamp = fmt.Sprintf("%d", now.UnixNano())

					var km KeyManager.KeyManager
					km.Init()
					km.GenRandomKeyPair()
					s.Address = km.GetAddress()
					s.Pubkey = km.GetPubkey()
					s.Sig = ""
					temp2, _ := json.Marshal(s)
					s.Sig, _= km.Sign(KeyManager.GetHash(temp2))
					mutex.Lock()
					messages = append(messages, s)
					mutex.Unlock()
				}
			}()
		}
		wg.Wait()
		wg1 := sync.WaitGroup{}
		wg1.Add(cs)


		url0 := "http://10.0.0.19:8010/CreateAccount"
		url1 := "http://10.0.0.20:8010/CreateAccount"
		url2 := "http://10.0.0.21:8010/CreateAccount"
		contentType := "application/json;charset=utf-8"
		client := &http.Client{
			Transport: &http.Transport{
				Dial: PrintLocalDial,
			},
		}
		for i := 0; i < cs; i++ {

			go func() {
				defer wg1.Done()

				for j := 0; j < gs; j++ {
					now := atomic.AddInt32(&num, 1)
					data, err := json.Marshal(messages[now-1])
					if err != nil {
						fmt.Println(err)
					}
					body0 := bytes.NewBuffer(data)
					body1 := bytes.NewBuffer(data)
					body2 := bytes.NewBuffer(data)
					resp, err :=client.Post(url0, contentType, body0)
					resp1, err := http.Post(url1, contentType, body1)
					resp2, err := http.Post(url2, contentType, body2)

					if err != nil {
						log.Println("Post failed:", err)
						return
					}
					_, _ = ioutil.ReadAll(resp.Body)
					_, _ = ioutil.ReadAll(resp1.Body)
					_, _ = ioutil.ReadAll(resp2.Body)
					defer resp.Body.Close()
					defer resp1.Body.Close()
					defer resp2.Body.Close()

					//data, _ := json.Marshal(basic)
					//fmt.Println(string(data))
				}
			}()
		}
		wg1.Wait()
		time.Sleep(700 * time.Millisecond)
		wg2 := sync.WaitGroup{}
		wg2.Add(cs)
		for i := 0; i < cs; i++ {

			go func() {
				defer wg2.Done()
				for j := 0; j < gs; j++ {
					now := atomic.AddInt32(&num, 1)
					data, err := json.Marshal(messages[now-1])
					if err != nil {
						fmt.Println(err)
					}
					body0 := bytes.NewBuffer(data)
					body1 := bytes.NewBuffer(data)
					body2 := bytes.NewBuffer(data)
					resp, err := http.Post(url0, contentType, body0)
					resp1, err := http.Post(url1, contentType, body1)
					resp2, err := http.Post(url2, contentType, body2)
					if err != nil {
						log.Println("Post failed:", err)
						return
					}
					_, _ = ioutil.ReadAll(resp.Body)
					_, _ = ioutil.ReadAll(resp1.Body)
					_, _ = ioutil.ReadAll(resp2.Body)
					defer resp.Body.Close()
					defer resp1.Body.Close()
					defer resp2.Body.Close()

					//data, _ := json.Marshal(basic)
					//fmt.Println(string(data))
				}
			}()
		}
		wg2.Wait()
		now := atomic.AddInt32(&num, 1)
		fmt.Println(now-1, "个交易发送完成")
	}
}

func PrintLocalDial(network, addr string) (net.Conn, error) {
	dial := net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	conn, err := dial.Dial(network, addr)
	if err != nil {
		return conn, err
	}

	fmt.Println("connect done, use", conn.LocalAddr().String())

	return conn, err
}
