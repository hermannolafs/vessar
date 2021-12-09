package samskipti

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	samskipti "hermannolafs/vessar/samskipti/go"
	"net"
	"strconv"
)

func Run(líffæri samskipti.HjartaServer, port int) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	samskipti.RegisterHjartaServer(srv, líffæri)
	reflection.Register(srv)

	if err := srv.Serve(listener); err != nil {
		panic(err)
	}
}
