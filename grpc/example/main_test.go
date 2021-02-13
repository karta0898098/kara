package example

import (
	"context"
	"testing"

	"github.com/karta0898098/kara/grpc"
	pb "github.com/karta0898098/kara/grpc/example/echo"
	"github.com/karta0898098/kara/tracer"
	"github.com/karta0898098/kara/zlog"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"

	"go.uber.org/fx"
)

type testSuite struct {
	suite.Suite
	app *fx.App
}

func TestEndpoint(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (s *testSuite) SetupTest() {
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

	s.app = fx.New(
		fx.Supply(config),
		fx.Supply(logConfig),
		fx.Provide(NewHandler),
		fx.Provide(grpc.NewGRPC),
		fx.Invoke(zlog.New),
		fx.Invoke(grpc.RunGRPC),
		fx.Invoke(SetGRPCService),
	)

	go s.app.Run()
}

func (s *testSuite) TearDownTest() {
	s.app.Done()
}

func (s *testSuite) TestHandler_Echo() {
	type args struct {
		ctx     context.Context
		req     *pb.EchoRequest
		traceID string
	}

	tests := []struct {
		name    string
		args    args
		want    *pb.EchoReply
		traceID string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: &pb.EchoRequest{
					Msg: "hello",
				},
				traceID: "12345678",
			},
			want: &pb.EchoReply{
				Msg: "hello",
			},
			traceID: "12345678",
		},
	}
	cc, _ := grpc.NewClient("127.0.0.1:8080")
	client := pb.NewEchoClient(cc)

	for _, tt := range tests {
		tt.args.ctx = context.WithValue(
			tt.args.ctx,
			tracer.TraceIDKey,
			tt.args.traceID,
		)
		reply, err := client.Echo(tt.args.ctx, tt.args.req)
		s.Equal(nil, err)
		s.Equal(tt.want.Msg, reply.Msg)
		// s.Equal(tt.traceID, trace.GetTraceID(tt.args.ctx))
	}
}
