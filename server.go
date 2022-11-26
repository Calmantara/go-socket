package main

import (
	"log"
	"net"
)

var (
	connMap = make(map[net.Conn][]byte)
)

func echoServer(listener net.Listener) {
	for {
		fd, err := listener.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		log.Printf("accepted file desc: %x\n", fd)
		go sendRecv(fd)
	}
}

func sendRecv(c net.Conn) {
	for {
		// reading data
		buf := make([]byte, 512)
		_, err := c.Read(buf)
		if err != nil {
			return
		}
		// if only \n received
		log.Printf("Server got:%v from %v\n", string(buf), c)
		for _, val := range buf {
			if val == 10 {
				if len(connMap[c]) >= 1 {
					data := connMap[c]
					if err := proceedResponse(c, data); err != nil {
						c.Close()
						delete(connMap, c)
						return
					}
				}
				connMap[c] = []byte{}
				continue
			}
			if val != 0 {
				connMap[c] = append(connMap[c], val)
			}
		}
	}
}

func proceedResponse(c net.Conn, data []byte) error {
	// validate and create response
	log.Printf("validating data:%s\n", data)
	req, err := validateRequest(data)
	if err != nil {
		log.Println("error validating request")
		return err
	}

	// sending callback data
	res, err := transformResponse(req)
	if err != nil {
		log.Println("error transforming response")
		return err
	}
	log.Printf("sending callback data:%s\n", res)

	_, err = c.Write(res)
	if err != nil {
		log.Fatal("Write: ", err)
	}
	return nil
}
