package main

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"

	dbConfig "github.com/dioxine/grpc-pb/internal/db"
	interfaces "github.com/dioxine/grpc-pb/pkg/v1"
	grpchandler "github.com/dioxine/grpc-pb/pkg/v1/handler/grpc"
	repo "github.com/dioxine/grpc-pb/pkg/v1/handler/repository"
	usecase "github.com/dioxine/grpc-pb/pkg/v1/handler/usecase"
	"github.com/pocketbase/pocketbase/daos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func initUserServer(db *daos.Dao) interfaces.UseCaseInterface {
	userRepo := repo.New(db)
	return usecase.New(userRepo)
}

func main() {
	db := dbConfig.PbInit()

	// section of ssl security

	cert, err := tls.LoadX509KeyPair("cert/server/public/server.crt", "cert/server/private/server.pem")
	if err != nil {
		log.Fatalf("failed to load key pair: %s", err)
	}

	ca := x509.NewCertPool()
	caFilePath := "cert/ca/public/ca.crt"
	caBytes, err := os.ReadFile(caFilePath)
	if err != nil {
		log.Fatalf("failed to read ca cert %q: %v", caFilePath, err)
	}
	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		log.Fatalf("failed to parse %q", caFilePath)
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    ca,
	}

	// start the grpc server

	grpcServer := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("ERROR STARTING THE SERVER : %v", err)
	}

	// attaching database and handlers

	userUseCase := initUserServer(db.Dao())
	grpchandler.NewServer(grpcServer, userUseCase)

	// start GRPC serving to the address as goroutine
	go grpcServer.Serve(lis)

	// start Pocketbase at main thread
	dbConfig.PbStart(db)

}
