package ports

import (
	pb "choirulanwar/user-svc/internal/core/domain"
)

type UserRepository interface {
	Find(id string) (*pb.FindRes, error)
	Store(data *pb.User) error
	Update(id string, data *pb.User) error
	FindAll(page int64, limit int64, orderBy string, orderType string) (*pb.FindAllRes, error)
	Delete(id string) error
}

type UserService interface {
	Find(id string) (*pb.FindRes, error)
	Store(data *pb.User) error
	Update(id string, data *pb.User) error
	FindAll(page int64, limit int64, orderBy string, orderType string) (*pb.FindAllRes, error)
	Delete(id string) error
}
