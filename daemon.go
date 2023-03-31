package gin_daemon

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// Daemon process
type Daemon struct {
	once       sync.Once
	unixSocket *os.File
}

func NewDaemon() *Daemon {
	d := &Daemon{}
	return d
}

func (d *Daemon) RunAsMaster() error {
	//TODO listen tcp http for outer
	//TODO listen unix http for internal
	unixSocket := filepath.Join(os.TempDir(), "gin-daemon-unix-http.sock")
	//defer os.Remove(unixSocket)

	r := gin.Default()
	r.GET("/health_check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	r.GET("/stop", func(c *gin.Context) {
		// TODO
	})

	r.GET("/run", func(c *gin.Context) {
		go d.RunPlugin()
	})

	os.Remove(unixSocket)
	ua, err := net.ResolveUnixAddr("unix", unixSocket)
	if err != nil {
		return err
	}

	l, err := net.ListenUnix("unix", ua)
	if err != nil {
		return err
	}
	file, err := l.File()
	if err != nil {
		return err
	}
	d.unixSocket = file
	http.Serve(l, r)

	return nil
}

func (d *Daemon) StopWorker() {

}

func (d *Daemon) RunPlugin() error {
	// TODO get metadata from master
	basePath := os.Args[0]
	cmd := exec.Command(basePath + "/plugin")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// pass unix server fd
	cmd.ExtraFiles = []*os.File{d.unixSocket}

	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func GinPlugin() gin.HandlerFunc {
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
	}
}
