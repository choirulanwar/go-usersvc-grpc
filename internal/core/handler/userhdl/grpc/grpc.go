package userhdlgrpc

import (
	// "github.com/micro/go-micro/v2/errors"

	"errors"
	"fmt"
	"os"

	grpcErrors "go-micro.dev/v4/errors"

	pb "choirulanwar/user-svc/internal/core/domain"
	"choirulanwar/user-svc/internal/core/ports"
	"choirulanwar/user-svc/pkg/constants"
	"context"
)

type GrpcHandler struct {
	userService ports.UserService
}

// Create a new GRPC handler
func NewGrpcHandler(userService ports.UserService) *GrpcHandler {
	return &GrpcHandler{
		userService: userService,
	}
}

// Find user by ID
func (hdl *GrpcHandler) Find(ctx context.Context, req *pb.FindReq, rsp *pb.FindRes) error {
	user, err := hdl.userService.Find(req.Id)

	if err != nil {
		return grpcErrors.NotFound(fmt.Sprintf("%s.%s", os.Getenv("SVC_NAME"), "find"), err.Error())
	}

	rsp.Data = user.Data

	return nil
}

// Store new user
func (hdl *GrpcHandler) Store(ctx context.Context, req *pb.StoreReq, rsp *pb.StoreRes) error {
	err := hdl.userService.Store(req.Data)

	if err != nil {
		if errors.Is(err, constants.ErrExists) {
			return grpcErrors.Conflict(fmt.Sprintf("%s.%s", os.Getenv("SVC_NAME"), "store"), err.Error())
		}

		return grpcErrors.BadRequest(fmt.Sprintf("%s.%s", os.Getenv("SVC_NAME"), "store"), err.Error())
	}

	return nil
}

// Update user
func (hdl *GrpcHandler) Update(ctx context.Context, req *pb.UpdateReq, rsp *pb.UpdateRes) error {
	err := hdl.userService.Update(req.Id, req.Data)

	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return grpcErrors.NotFound(fmt.Sprintf("%s.%s", os.Getenv("SVC_NAME"), "update"), err.Error())
		}

		return grpcErrors.BadRequest(fmt.Sprintf("%s.%s", os.Getenv("SVC_NAME"), "update"), err.Error())
	}

	return nil
}

// Find all users
func (hdl *GrpcHandler) FindAll(ctx context.Context, req *pb.FindAllReq, rsp *pb.FindAllRes) error {
	users, err := hdl.userService.FindAll(req.Page, req.Limit, req.OrderBy.String(), req.OrderType.String())

	if err != nil {
		return grpcErrors.BadRequest(fmt.Sprintf("%s.%s", os.Getenv("SVC_NAME"), "findAll"), err.Error())
	}

	rsp.TotalDatas = users.TotalDatas
	rsp.Limit = users.Limit
	rsp.Page = users.Page
	rsp.TotalPages = users.TotalPages
	rsp.Datas = users.Datas
	rsp.NextPage = users.NextPage
	rsp.PrevPage = users.PrevPage

	return nil
}

// Delete user
func (hdl *GrpcHandler) Delete(ctx context.Context, req *pb.DeleteReq, rsp *pb.DeleteRes) error {
	err := hdl.userService.Delete(req.Id)

	if err != nil {
		return grpcErrors.BadRequest(fmt.Sprintf("%s.%s", os.Getenv("SVC_NAME"), "delete"), err.Error())
	}

	return nil
}
