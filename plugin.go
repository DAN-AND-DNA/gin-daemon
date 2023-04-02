package gin_daemon

import (
	"context"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
)

type Plugin struct {
	unixHttpServer *http.Server
	l              net.Listener
}

func NewPlugin() *Plugin {
	p := &Plugin{}
	f := os.NewFile(uintptr(3), "")
	l, err := net.FileListener(f)
	if err != nil {
		panic(err)
	}

	p.l = l

	r := gin.Default()
	r.Use(RunAsPlugin(p))
	unixHttpServer := &http.Server{
		Handler: r,
	}
	p.unixHttpServer = unixHttpServer
	return p
}

func (p *Plugin) Run() error {
	err := p.unixHttpServer.Serve(p.l)
	if err != nil {
		return err
	}

	return nil
}

func (p *Plugin) Stop() {
	p.unixHttpServer.Shutdown(context.TODO())
}

func RunAsPlugin(plugin *Plugin) gin.HandlerFunc {
	return func(c *gin.Context) {
		/* TODO
		xx/stop 停止父进程，子进程的心跳根据父进程的pid判断，是否要优雅关闭
		xx/reload 新创建一个子进程，父进程把请求打给新进程，老的子进程优雅关闭

		父进程监听http请求
		子进程监听http请求

		钩子，当子进程成功创建，父进程执行握手的钩子
		unixTestSocket := filepath.Join(os.TempDir(), "unix_unit_test")
		defer os.Remove(unixTestSocket)

		r := New()
		r.RunUnix(unixTestSocket)

		c, err := net.Dial("unix", unixTestSocket)
		*/
		fullPath := c.FullPath()
		switch fullPath {
		case "/health_check":
			c.JSON(200, gin.H{
				"message": "ok",
			})

			return
		case "/stop":
			plugin.Stop()
		}

	}
}
