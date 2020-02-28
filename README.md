### 使用

```golang
func main() {
	listen := &tcp.TcpServer {
		Addr : "0.0.0.0",
		Port : "8888",
		NewConn : func(ctx *ConnContext) {
			fmt.Println(ctx.Conn)
		}
		NewMessage : func(ctx *ConnContext,message string) {
			fmt.Println(message)
		}
		ConnClose : func(ctx *ConnContext,err error) {

		}
	}

	listen.Listen()
}


```