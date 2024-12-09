package main

import (
	psqlv1 "github.com/crewlinker/protosql/psql/v1"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// queryOpts returns our plugin specific response options.
func queryOpts(m *protogen.Message) *psqlv1.QueryOptions {
	opts, hasOpts := m.Desc.Options().(*descriptorpb.MessageOptions)
	if !hasOpts {
		return nil
	}

	ext, hasOpts := proto.GetExtension(opts, psqlv1.E_Query).(*psqlv1.QueryOptions)
	if !hasOpts {
		return nil
	}

	return ext
}
