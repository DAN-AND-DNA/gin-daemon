package gin_daemon

import "github.com/gin-gonic/gin"

func GinDaemon() gin.HandlerFunc {
	return func(c *gin.Context) {
		/* TODO
		xx/run 启动子进程
		xx/stop 停止父进程，子进程的心跳根据父进程的pid判断，是否要优雅关闭
		xx/reload 新创建一个子进程，父进程把流量打给新进程，老的子进程优雅关闭

		父进程监听http请求
		子进程监听http请求

		钩子，当子进程成功创建，父进程执行握手的钩子

		*/
	}
}
