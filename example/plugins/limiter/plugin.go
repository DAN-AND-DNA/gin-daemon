package limiter

import gin_daemon "github.com/dan-and-dna/gin-daemon"

type Limiter struct {
}

func main() {
	p := gin_daemon.NewPlugin()
	err := p.Run()
	if err != nil {
		panic(err)
	}
}
