package share

import (
	"avrilko-rpc/codec"
	"avrilko-rpc/protocol"
)

const (
	AuthKey = "__AUTH"
)

type ContextKey string

var ReqMetaDataKey = ContextKey("__req_metadata")

var ResMetaDataKey = ContextKey("__res_metadata")

var (
	Codecs = map[protocol.SerializeType]codec.Codec{
		protocol.SerializeNone: &codec.ByteCodec{},
		protocol.JSON:          &codec.JSONCodec{},
		protocol.ProtoBuffer:   &codec.PBCodec{},
		protocol.MsgPack:       &codec.MsgpackCodec{},
	}
)
