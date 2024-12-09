package main

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bufbuild/protoplugin"
	"github.com/iancoleman/strcase"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	psqlv1 "github.com/crewlinker/protosql/psql/v1"
)

func main() {
	protoplugin.Main(protoplugin.HandlerFunc(handle))
}

func handle(
	_ context.Context,
	_ protoplugin.PluginEnv,
	resp protoplugin.ResponseWriter,
	req protoplugin.Request,
) error {
	resp.SetFeatureProto3Optional()
	resp.SetFeatureSupportsEditions(descriptorpb.Edition_EDITION_2023, descriptorpb.Edition_EDITION_2024)

	fileDescs, err := req.FileDescriptorsToGenerate()
	if err != nil {
		return err
	}

	logs := slog.New(slog.NewTextHandler(os.Stderr, nil))

	for _, fileDesc := range fileDescs {
		messages := fileDesc.Messages()
		for idx := range messages.Len() {
			msgDesc := messages.Get(idx)

			options := msgDesc.Options()
			if options == nil {
				continue
			}

			ext := proto.GetExtension(options, psqlv1.E_Query)
			if ext == nil {
				logs.Warn("Error retrieving extension", "message", msgDesc.Name(), "error", err)
				continue
			}

			if queryOptions, ok := ext.(*psqlv1.QueryOptions); ok {
				switch queryOptions.GetKind() {
				case psqlv1.QueryKind_QUERY_KIND_EXEC:
					fpath := filepath.Join(filepath.Dir(fileDesc.Path()), strcase.ToSnake(string(msgDesc.Name()))+".sql")

					logs.Info("exec", slog.String("path", fpath))
				case psqlv1.QueryKind_QUERY_KIND_QUERY:
					logs.Info("query")
				case psqlv1.QueryKind_QUERY_KIND_QUERY_ROW:
					logs.Info("query row")
				default:
					continue
				}
			} else {
				logs.Warn("Invalid QueryOptions type", "message", msgDesc.Name())
				continue
			}
		}
	}

	return nil
}
