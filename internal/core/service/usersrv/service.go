package usersrv

import (
	"time"

	// "github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	pb "choirulanwar/user-svc/internal/core/domain"
	"choirulanwar/user-svc/internal/core/ports"
	"choirulanwar/user-svc/pkg/encrypt"
	"choirulanwar/user-svc/pkg/uidgen"
)

type service struct {
	userRepo ports.UserRepository
	uidGen   uidgen.UIDGen
	hash     encrypt.Hash
}

// Create service
func NewUserService(userRepo ports.UserRepository, uidGen uidgen.UIDGen) ports.UserService {
	return &service{userRepo: userRepo, uidGen: uidGen}
}

// Find user by id
func (s *service) Find(id string) (*pb.FindRes, error) {
	return s.userRepo.Find(id)
}

// Store new user
func (s *service) Store(data *pb.User) error {
	hashedPassword, err := s.hash.Generate(data.Password)
	if err != nil {
		return errors.New(err.Error())
	}

	data.Id = s.uidGen.New()
	data.Password = hashedPassword
	data.Role = 1
	data.IsEmailVerified = false
	data.IsActive = false
	data.ActiveUntil = time.Now().Unix()
	data.CreatedAt = time.Now().Unix()
	data.UpdatedAt = time.Now().Unix()

	// validate := validator.New()
	// notValid := validate.Struct(user)
	// if notValid != nil {
	// 	return errors.New(notValid.Error())
	// }

	return s.userRepo.Store(data)
}

// Update user
func (s *service) Update(id string, data *pb.User) error {
	data.UpdatedAt = time.Now().Unix()

	return s.userRepo.Update(id, data)
}

// Find all user
func (s *service) FindAll(offset int64, limit int64, orderBy string, orderType string) (*pb.FindAllRes, error) {
	return s.userRepo.FindAll(offset, limit, orderBy, orderType)
}

// Delete user
func (s *service) Delete(id string) error {
	return s.userRepo.Delete(id)
}
