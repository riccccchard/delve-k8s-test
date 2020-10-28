package utils

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func Pb2Json(pbMessage proto.Message) (string, error) {
	marshaler := jsonpb.Marshaler{}
	marshaler.EmitDefaults = true //这个设置最关键
	jsons, err := marshaler.MarshalToString(pbMessage)
	return jsons, err
}
