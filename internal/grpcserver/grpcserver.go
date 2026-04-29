package grpcserver

import (
	"context"

	"configlinter/internal/engine"
	"configlinter/internal/parser"
	pb "configlinter/proto/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedConfigLinterServer
	registry *parser.Registry
	engine   *engine.Engine
}

func New(reg *parser.Registry, eng *engine.Engine) *Server {
	return &Server{
		registry: reg,
		engine:   eng,
	}
}

func (s *Server) CheckConfig(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	if req.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "empty content")
	}

	format := protoFormatToString(req.Format)
	if format == "" {
		return nil, status.Error(codes.InvalidArgument, "unsupported format")
	}

	p, err := s.registry.GetByFormat(format)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "unsupported format: %s", format)
	}

	root, err := p.Parse([]byte(req.Content))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "parse error: %s", err)
	}

	findings := s.engine.Analyze(root)

	resp := &pb.CheckResponse{}
	for _, f := range findings {
		resp.Issues = append(resp.Issues, &pb.Issue{
			RuleId:         f.RuleID,
			Severity:       f.Severity.ToString(),
			Path:           f.Path,
			Message:        f.Message,
			Recommendation: f.Recomendation,
		})
	}

	return resp, nil
}

func protoFormatToString(f pb.Format) string {
	switch f {
	case pb.Format_FORMAT_JSON:
		return "json"
	case pb.Format_FORMAT_YAML:
		return "yaml"
	case pb.Format_FORMAT_TOML:
		return "toml"
	default:
		return ""
	}
}
