package gin_daemon

const (
	PluginError MsgType = iota + 1001
	PluginHeartBeat
)

type MsgType int

type Msg struct {
	Type MsgType
	Err  error
}
