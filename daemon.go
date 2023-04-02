package gin_daemon

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

// Daemon process
type Daemon struct {
	lUnix         *net.UnixListener
	tcpHttpServer *http.Server
	wg            sync.WaitGroup
	MsgChan       chan Msg
	Pid           int
	PluginPid     int
}

func NewDaemon() *Daemon {
	d := &Daemon{
		MsgChan: make(chan Msg, 100),
		Pid:     os.Getpid(),
	}

	d.Poll()
	lUnix, err := newUnixListener()
	if err != nil {
		panic(err)
	}
	d.lUnix = lUnix

	r := gin.Default()
	r.Use(runAsDaemon(d))
	tcpHttpServer := &http.Server{
		Addr:    ":3777",
		Handler: r,
	}
	d.tcpHttpServer = tcpHttpServer

	return d
}

func (d *Daemon) Run() error {
	err := d.tcpHttpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (d *Daemon) Poll() {
	go func() {
		d.wg.Add(1)
		defer d.wg.Done()

		t1 := time.Tick(1 * time.Second)
		t3 := time.Tick(3 * time.Second)

		for {
			select {
			case msg := <-d.MsgChan:
				switch msg.Type {
				case PluginError:
				case PluginHeartBeat:
				}
			case <-t1:
			case <-t3:
				d.watchPlugin()
			}
		}
	}()
}

func (d *Daemon) watchPlugin() {
	if d.IsPluginRunning() {
		return
	}

	err := d.RunPlugin()
	if err != nil {
		log.Println(err)
		return
	}
}

func newUnixListener() (*net.UnixListener, error) {
	unixSocket := filepath.Join(os.TempDir(), "plugin.sock")
	_ = os.Remove(unixSocket)

	ua, err := net.ResolveUnixAddr("unix", unixSocket)
	if err != nil {
		return nil, err
	}

	l, err := net.ListenUnix("unix", ua)
	if err != nil {
		return l, err
	}

	return l, nil
}

func (d *Daemon) IsPluginRunning() bool {
	p, err := os.FindProcess(d.PluginPid)
	if err != nil {
		err = p.Signal(syscall.Signal(0))
		return err != nil
	}

	return false
}

// StopPlugin shutdown plugin process gracefully
func (d *Daemon) StopPlugin() error {
	return nil
}

// RunPlugin run new plugin process
func (d *Daemon) RunPlugin() error {
	lFile, err := d.lUnix.File()
	if err != nil {
		return err
	}

	basePath := os.Args[0]
	cmd := exec.Command(basePath + "/plugin")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{lFile}

	err = cmd.Start()
	if err != nil {
		return err
	}

	d.wg.Add(1)
	defer d.wg.Done()

	go func() {
		d.wg.Add(1)
		defer d.wg.Done()

		err = cmd.Wait()
		if err != nil {
			d.MsgChan <- Msg{Type: PluginError, Err: err}
			return
		}
	}()

	return nil
}

func runAsDaemon(d *Daemon) gin.HandlerFunc {
	return func(c *gin.Context) {
		fullPath := c.FullPath()
		switch fullPath {
		case "/stop":
			// TODO
			d.tcpHttpServer.Shutdown(context.TODO())

			d.wg.Wait()

		case "/reload":
			err := d.StopPlugin()
			if err != nil {
				log.Println(err)
			}

			err = d.RunPlugin()
			if err != nil {
				log.Println(err)
			}
		}
	}
}
