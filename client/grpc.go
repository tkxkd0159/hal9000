package client

import (
	"crypto/x509"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func OpengRPCConn(addr string, tls bool) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error

	if tls {
		var pool *x509.CertPool
		pool, err = x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
		creds := credentials.NewClientTLSFromCert(pool, "")

		conn, err = grpc.Dial(
			addr,
			grpc.WithTransportCredentials(creds),
		)
	} else {
		conn, err = grpc.Dial(
			addr,
			grpc.WithInsecure(),
		)
	}

	return conn, err
}
