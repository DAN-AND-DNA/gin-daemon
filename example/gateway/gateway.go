package gateway

import (
	gin_daemon "github.com/dan-and-dna/gin-daemon"
	"github.com/dan-and-dna/gin-daemon/example/pb"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
)

// Gateway 连接管理、数据存储和消息转发
type Gateway struct {
}

func NewGateway() *Gateway {
	gw := &Gateway{}
	return gw
}

func (gw *Gateway) Run() error {
	l, err := net.Listen("tcp", "0.0.0.0:7001")
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		process(conn)
	}

	return nil
}

func process(conn net.Conn) {
	if conn == nil {
		return
	}
	defer conn.Close()

	var cache = make([]byte, 2048)
	var buf = make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		data := buf[:n]
		_ = data

		msg := &pb.Student{}
		proto.Unmarshal(cache, msg)
	}
}

func main() {
	daemon := gin_daemon.NewDaemon()
	err := daemon.Run()
	if err != nil {
		panic(err)
	}
}
