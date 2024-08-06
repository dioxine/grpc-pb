package grpc_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"testing"

	pb "github.com/dioxine/grpc-pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var id string = "wcqsspl3d3obmh9"

func Connection(t *testing.T) *grpc.ClientConn {
	// section of ssl security

	cert, err := tls.LoadX509KeyPair("../../../../../cert/client/public/client.crt", "../../../../../cert/client/private/client.pem")
	if err != nil {
		log.Fatalf("failed to load client cert: %v", err)
	}

	ca := x509.NewCertPool()
	caFilePath := "../../../../../cert/ca/public/ca.crt"
	caBytes, err := os.ReadFile(caFilePath)
	if err != nil {
		log.Fatalf("failed to read ca cert %q: %v", caFilePath, err)
	}
	if ok := ca.AppendCertsFromPEM(caBytes); !ok {
		log.Fatalf("failed to parse %q", caFilePath)
	}

	tlsConfig := &tls.Config{
		ServerName:   "localhost",
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
	}

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		t.Fatal("the connection with the server cannot be established")
	}
	return conn
}

func TestCreateUser(t *testing.T) {
	conn := Connection(t)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	createRequest := &pb.CreateUserRequest{
		Username: "pjburger",
		Name:     "Pj Burger",
		Email:    "pjburger@mail.ru",
		Password: "123456789",
	}

	createResponse, err := client.Create(context.Background(), createRequest)
	if err != nil {
		t.Fatalf("CREATE FAILED: %v", err)
	}

	t.Log("CREATE SUCCESS: ", createResponse)

	id = createResponse.GetId()
}

func TestReadUser(t *testing.T) {
	conn := Connection(t)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	readRequest := &pb.SingleUserRequest{
		Id: id,
	}

	readResponse, err := client.Read(context.Background(), readRequest)
	fmt.Println(readResponse)
	if err != nil {
		t.Fatalf("READ FAILED: %v", err)
	}

	t.Log("READ SUCCESS: ", readResponse)
}

func TestUpdateUser(t *testing.T) {
	conn := Connection(t)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	updateRequest := &pb.UpdateUserRequest{
		Id:       id,
		Username: "TestUser2",
		Name:     "Test User2",
		Email:    "test_user@mail.ru",
		Password: "987654321",
	}

	updateResponse, err := client.Update(context.Background(), updateRequest)

	if err != nil {
		t.Fatalf("UPDATE FAILED: %v", err)
	}

	t.Log("UPDATE SUCCESS: ", updateResponse)
}

func TestDeleteUser(t *testing.T) {
	conn := Connection(t)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	deleteRequest := &pb.SingleUserRequest{
		Id: id,
	}

	deleteResponse, err := client.Delete(context.Background(), deleteRequest)

	if err != nil {
		t.Fatalf("DELETE FAILED: %v", err)
	}

	t.Log("DELETE SUCCESS: ", deleteResponse)
}
