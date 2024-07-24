package grpc

import (
	"context"
	"errors"

	"github.com/dioxine/grpc-pb/internal/models"
	interfaces "github.com/dioxine/grpc-pb/pkg/v1"
	pb "github.com/dioxine/grpc-pb/proto"
	"google.golang.org/grpc"
)

type UserServStruct struct {
	useCase interfaces.UseCaseInterface
	pb.UnimplementedUserServiceServer
}

func NewServer(grpcServer *grpc.Server, usecase interfaces.UseCaseInterface) {
	userGrpc := &UserServStruct{useCase: usecase}
	pb.RegisterUserServiceServer(grpcServer, userGrpc)
}

// Create
//
// This function creates a user whose data is supplied
// through the CreateUserRequest message of proto
func (srv *UserServStruct) Create(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserProfileResponse, error) {

	data := srv.transformUserRPC(req)
	if data.Email == "" || data.Name == "" {
		return &pb.UserProfileResponse{}, errors.New("please provide all fields")
	}

	user, err := srv.useCase.Create(data)
	if err != nil {
		return &pb.UserProfileResponse{}, err
	}
	return srv.transformUserModel(user), nil
}

// Read
//
// This function returns the user instance of which ID
// is supplied through the SingleUserRequest message field of proto
func (srv *UserServStruct) Read(ctx context.Context, req *pb.SingleUserRequest) (*pb.UserProfileResponse, error) {
	id := req.GetId()
	if id == "" {
		return &pb.UserProfileResponse{}, errors.New("id cannot be blank")
	}
	user, err := srv.useCase.Get(id)
	if err != nil {
		return &pb.UserProfileResponse{}, err
	}
	return srv.transformUserModel(user), nil
}

// Update
//
// This function returns the success message if user is updated with
// data supplied through the UpdateUserRequest message of proto
func (srv *UserServStruct) Update(ctx context.Context, req *pb.UpdateUserRequest) (*pb.SuccessResponse, error) {

	data := srv.transformUpdateUserRPC(req)

	if data.Id == "" {
		return &pb.SuccessResponse{}, errors.New("id cannot be blank")
	}

	if data.Email == "" || data.Name == "" || data.Username == "" || data.Password == "" {
		return &pb.SuccessResponse{}, errors.New("please provide all fields")
	}

	err := srv.useCase.Update(data)

	if err != nil {
		return &pb.SuccessResponse{}, errors.New("email cannot be changed")
	} else {
		return &pb.SuccessResponse{Response: "successfully updated user"}, nil
	}
}

// Delete
//
// This function returns the success message if user is deleted
func (srv *UserServStruct) Delete(ctx context.Context, req *pb.SingleUserRequest) (*pb.SuccessResponse, error) {
	id := req.GetId()
	if id == "" {
		return &pb.SuccessResponse{}, errors.New("id cannot be blank")
	}
	err := srv.useCase.Delete(id)
	if err != nil {
		return &pb.SuccessResponse{}, err
	}
	return &pb.SuccessResponse{Response: "successfully deleted user"}, nil
}

// Helper functions to transform requests and responses
func (srv *UserServStruct) transformUserRPC(req *pb.CreateUserRequest) models.User {
	return models.User{Username: req.GetUsername(), Name: req.GetName(), Email: req.GetEmail(), Password: req.GetPassword()}
}

func (srv *UserServStruct) transformUpdateUserRPC(req *pb.UpdateUserRequest) models.User {
	return models.User{Id: req.GetId(), Username: req.GetUsername(), Name: req.GetName(), Email: req.GetEmail(), Password: req.GetPassword()}
}
func (srv *UserServStruct) transformUserModel(user models.User) *pb.UserProfileResponse {
	return &pb.UserProfileResponse{Id: string(user.Id), Name: user.Name, Email: user.Email}
}
