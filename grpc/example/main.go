package example

import (
	"context"
	"time"

	"github.com/karta0898098/kara/grpc"
	pb "github.com/karta0898098/kara/grpc/example/echo"
	"github.com/karta0898098/kara/zlog"
	"github.com/rs/zerolog"

	"go.uber.org/fx"
	rpc "google.golang.org/grpc"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (e *Handler) Echo(ctx context.Context, req *pb.EchoRequest) (resp *pb.EchoReply, err error) {
	return &pb.EchoReply{
		Msg:      req.Msg,
		Unixtime: time.Now().Unix(),
	}, nil
}

func main() {
	config := &grpc.Config{
		Mode:        "debug",
		Port:        ":8080",
		RequestDump: true,
	}

	logConfig := &zlog.Config{
		Env:   "local",
		AppID: "test-grpc",
		Level: int8(zerolog.DebugLevel),
		Debug: true,
	}

	app := fx.New(
		fx.Supply(config),
		fx.Supply(logConfig),
		fx.Provide(NewHandler),
		fx.Provide(grpc.NewGRPC),
		fx.Invoke(zlog.Setup),
		fx.Invoke(grpc.RunGRPC),
		fx.Invoke(SetGRPCService),
	)

	app.Run()
}

// SetGRPCService register gRPC handler
func SetGRPCService(server *rpc.Server, handler *Handler) {
	pb.RegisterEchoServer(server, handler)
}
