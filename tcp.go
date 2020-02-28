package handle

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

type (
	TcpServer struct {
		Addr  string
		Port  string
		NewConn func(ctx *ConnContext)
		NewMessage func(ctx *ConnContext,message string)
		ConnClose func(ctx *ConnContext,err error)
	}

	ConnContext struct {
		Conn net.Conn
		Server *TcpServer
		read   *bufio.Reader
		write  *bufio.Writer
	}

)

func (c *ConnContext) Send(message string) error {
	nn, err := c.write.Write([]byte(message))
	err = c.write.Flush()

	if err != nil {
		c.Conn.Close()
	}

	if nn < len([]byte(message)) {
		c.Conn.Close()
		return io.ErrShortWrite
	}

	return err
}

func (c *ConnContext) SendBytes(b []byte) error {
	nn, err := c.write.Write(b)
	err = c.write.Flush()
	if err != nil {
		c.Conn.Close()
	}

	if nn < len(b) {
		c.Conn.Close()
		return io.ErrShortWrite
	}

	return err
}

//listen tcp
func (t *TcpServer) Listen() {
	if t.Addr == "" || t.Port == "" {
		log.Fatal("Config not null.")
	}
	listener, err := net.Listen("tcp",fmt.Sprintf("%s:%s",t.Addr,t.Port))
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()
	log.Println("starting TCP server success.")
	for {
		conn, _ := listener.Accept()

		connContext := &ConnContext{
			Conn:conn,
			Server:t,
			read: bufio.NewReader(conn),
			write: bufio.NewWriter(conn),
		}


		go connContext.handle()
	}

}


func (c *ConnContext) handle() {
	defer c.Conn.Close()

	c.Server.NewConn(c)

	reader := json.NewDecoder(c.Conn)

	for {
		maps := make(map[string]interface{})
		err := reader.Decode(&maps)
		if err != nil {
			c.Server.ConnClose(c,err)
			c.Conn.Close()

			log.Println("conn read error")
			return
		}
		jso, _ := json.Marshal(maps)
		c.Server.NewMessage(c,string(jso))
	}
}