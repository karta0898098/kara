package logging

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var marshaler = jsonpb.Marshaler{
	EnumsAsInts:  false,
	EmitDefaults: true,
	OrigName:     true,
}

type protoMessage struct {
	msg proto.Message
}

func (m protoMessage) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if err := marshaler.Marshal(&buf, m.msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func requestDump(ctx context.Context, info *grpc.UnaryServerInfo, request bool, logger zerolog.Logger, msg interface{}, err error) {
	if !request {
		return
	}

	dict := zerolog.Dict()

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		header, err := json.Marshal(&md)
		if err == nil {
			dict.RawJSON("header", header)
			// logger = logger.With().RawJSON("header", header).Logger()
		}
	}

	protoMsg, ok := msg.(proto.Message)
	if ok {
		msg := protoMessage{protoMsg}
		buf, err := msg.MarshalJSON()
		if err == nil {
			dict.RawJSON("body", buf)
			logger.Info().Dict("dump", dict).Msg("grpc request dump.")
		}
	}
}

func replayDump(ctx context.Context, info *grpc.UnaryServerInfo, request bool, logger zerolog.Logger, msg interface{}, err error) {
	if !request {
		return
	}

	protoMsg, ok := msg.(proto.Message)
	if ok {
		msg := protoMessage{protoMsg}
		buf, _ := msg.MarshalJSON()
		logger.Info().Dict("dump", zerolog.Dict().RawJSON("body", buf)).Msg("grpc replay dump.")
	}
}
