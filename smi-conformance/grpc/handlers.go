package grpc

import (
	"context"
	"strconv"

	"github.com/layer5io/learn-layer5/smi-conformance/proto"
	test_gen "github.com/layer5io/learn-layer5/smi-conformance/test-gen"
	common "github.com/layer5io/meshkit/protobuf/common"
	controller "github.com/layer5io/meshkit/protobuf/controller"
)

var (
	maeshConfig = &test_gen.Maesh{
		PortSvcA: "9091",
		PortSvcB: "9091",
		PortSvcC: "9091",
	}

	linkerdConfig = &test_gen.Linkerd{
		PortSvcA: "9091",
		PortSvcB: "9091",
		PortSvcC: "9091",
	}
	istioConfig = &test_gen.Istio{
		PortSvcA: "9091",
		PortSvcB: "9091",
		PortSvcC: "9091",
	}
	osmConfig = &test_gen.OSM{
		PortSvcA: "9091",
		PortSvcB: "9091",
		PortSvcC: "9091",
	}
)

func (s *Service) RunTest(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	var config test_gen.ServiceMesh

	config = linkerdConfig
	switch req.Mesh.Name {
	case "linkerd":
		config = linkerdConfig
		req.Mesh.Annotations["linkerd.io/inject"] = "enabled"
	case "maesh":
		config = maeshConfig
	case "istio":
		config = istioConfig
		req.Mesh.Labels["istio-injection"] = "enabled"
	case "osm":
		config = osmConfig
		req.Mesh.Labels["openservicemesh.io/monitored-by"] = "osm"
	}

	result := test_gen.RunTest(config, req.Mesh.Annotations, req.Mesh.Labels)
	totalcases := 3
	failures := 0

	details := make([]*proto.Detail, 0)
	for _, res := range result.Testsuite[0].Testcase {
		d := &proto.Detail{
			Smispec:    res.Name,
			Duration:   res.Time,
			Assertions: strconv.Itoa(res.Assertions),
			Status:     "Passing",
		}
		if len(res.Failure.Text) > 2 {
			d.Result.Error = &common.CommonError{
				ShortDescription: res.Failure.Text,
			}
			d.Result.Message = res.Failure.Message
			d.Status = "Failing"
			failures += 1
		}
		details = append(details, d)
	}

	capability := "None"
	if totalcases-failures > totalcases/2 {
		capability = "Half"
	} else if failures == 0 {
		capability = "None"
	}

	return &proto.Response{
		Casespassed: strconv.Itoa(totalcases - failures),
		Passpercent: strconv.Itoa(((totalcases - failures) / totalcases) * 100),
		Mesh:        req.Mesh,
		Capability:  capability,
		Details:     details,
	}, nil
}

func (s *Service) Health(ctx context.Context) (*controller.ControllerHealth, error) {
	return &controller.ControllerHealth{}, nil
}

func (s *Service) Info(ctx context.Context) (*controller.ControllerInfo, error) {
	return &controller.ControllerInfo{}, nil
}
