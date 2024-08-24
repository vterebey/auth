package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/brianvoe/gofakeit"
	desc "github.com/vterebey/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedUserV1Server
}

// UserStore is responsible for storing user data.
type UserStore struct {
	users map[int64]*desc.User
}

// NewUserStore creates a new instance of UserStore.
func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[int64]*desc.User),
	}
}

// SaveUser saves a user to the store.
func (store *UserStore) SaveUser(user *desc.User) {
	store.users[user.Id] = user
}

// GetUser retrieves a user by ID.
func (store *UserStore) GetUser(ID int64) (user *desc.User) {
	return store.users[ID]
}

var userStore = NewUserStore()

func (s server) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Create User: %s", req.Info.Name)
	user := &desc.User{
		Id: int64(gofakeit.Number(0, 100)),
		Info: &desc.UserInfo{
			Name:  req.Info.Name,
			Email: req.Info.Email,
			Role:  req.Info.Role,
		},
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}
	userStore.SaveUser(user)
	return &desc.CreateResponse{Id: user.Id}, nil

}

func (s *server) Get(_ context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Printf("Get user with id: %d", req.Id)

	user := userStore.GetUser(req.GetId())
	return &desc.GetResponse{
		Info: &desc.User{
			Id: user.Id,
			Info: &desc.UserInfo{
				Name:  user.Info.Name,
				Email: user.Info.Email,
				Role:  user.Info.Role,
			},
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
