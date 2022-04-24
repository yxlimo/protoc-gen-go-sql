package main

import (
	"flag"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/yxlimo/protoc-gen-go-sql/internal/gen"
)

var (
	useEnumNumbers  = flag.Bool("use_enum_numbers", false, "render enums as integers as opposed to strings")
	emitUnpopulated = flag.Bool("emit_unpopulated", false, "render fields with zero values")
	userProtoNames  = flag.Bool("use_proto_names", true, "use original (.proto) name for fields")
	allowPartial    = flag.Bool("allow_partial", false, "allow partial results")
	discardUnknown  = flag.Bool("discard_unknown", true, "allow messages to contain unknown fields when unmarshaling")
)

func main() {
	pgs.Init(pgs.DebugEnv("DEBUG")).RegisterModule(gen.New(gen.Option{
		Marshaler: protojson.MarshalOptions{
			AllowPartial:    *allowPartial,
			UseProtoNames:   *userProtoNames,
			UseEnumNumbers:  *useEnumNumbers,
			EmitUnpopulated: *emitUnpopulated,
		},
		Unmarshaler: protojson.UnmarshalOptions{
			AllowPartial:   *allowPartial,
			DiscardUnknown: *discardUnknown,
		},
	})).RegisterPostProcessor(
		pgsgo.GoFmt(),
	).Render()
}
