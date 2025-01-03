package server

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/jwttoken"
	pb "github.com/MagicNetLab/ya-practicum-shortener/pkg/shortener_grpc"
)

// TestServer_UserRegister тест регистрации пользователя
func TestServer_UserRegister(t *testing.T) {
	s := Server{}
	err := config.Initialize()
	assert.NoError(t, err)
	conf := config.GetParams()
	err = repo.Initialize(conf)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	req := pb.UserRegisterRequest{Login: "login", Secret: "secret"}
	resp, err := s.UserRegister(ctx, &req)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Token)
	assert.True(t, jwttoken.ValidateToken(resp.Token, conf.GetJWTSecret()))

	resp, err = s.UserRegister(ctx, &req)
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.ErrorIs(t, err, status.Error(codes.AlreadyExists, "login already used"))
}

// TestServer_UserAuth тест авторизации пользователя
func TestServer_UserAuth(t *testing.T) {
	s := Server{}
	err := config.Initialize()
	assert.NoError(t, err)
	conf := config.GetParams()
	err = repo.Initialize(conf)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	ok, err := repo.CreateUser(ctx, "login", "secret")
	assert.NoError(t, err)
	assert.True(t, ok)

	req := pb.UserAuthRequest{Login: "login", Secret: "secret"}
	resp, err := s.UserAuth(ctx, &req)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Token)
	assert.True(t, jwttoken.ValidateToken(resp.Token, conf.GetJWTSecret()))

	req = pb.UserAuthRequest{Login: "login", Secret: "secret2"}
	resp, err = s.UserAuth(ctx, &req)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, status.Error(codes.NotFound, "failed authentication"))
}

// TestServer_EncodeLink тест сокращения одной ссылки
func TestServer_EncodeLink(t *testing.T) {
	s := Server{}
	err := config.Initialize()
	assert.NoError(t, err)
	conf := config.GetParams()
	err = repo.Initialize(conf)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ok, err := repo.CreateUser(ctx, "login", "secret")
	assert.NoError(t, err)
	assert.True(t, ok)

	userID, err := repo.AuthUser(ctx, "login", "secret")
	assert.NoError(t, err)
	token, err := jwttoken.GenerateToken(userID, conf.GetJWTSecret())
	assert.NoError(t, err)
	cancel()

	req := pb.EncodeLinkRequest{OriginalUrl: "https://yandex.ru"}
	header := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewIncomingContext(context.Background(), header)

	resp, err := s.EncodeLink(ctx, &req)
	assert.NoError(t, err)
	assert.NotNil(t, resp.ShortLink)
	assert.Contains(t, resp.ShortLink, conf.GetShortHost())

	resp, err = s.EncodeLink(ctx, &req)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, status.Error(codes.AlreadyExists, "link already exists"))
}

// TestServer_BatchEncodeLink тест сокращения пакета ссылок
func TestServer_BatchEncodeLink(t *testing.T) {
	s := Server{}
	err := config.Initialize()
	assert.NoError(t, err)
	conf := config.GetParams()
	err = repo.Initialize(conf)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ok, err := repo.CreateUser(ctx, "login", "secret")
	assert.NoError(t, err)
	assert.True(t, ok)

	userID, err := repo.AuthUser(ctx, "login", "secret")
	assert.NoError(t, err)
	token, err := jwttoken.GenerateToken(userID, conf.GetJWTSecret())
	assert.NoError(t, err)
	cancel()

	header := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewIncomingContext(context.Background(), header)

	req := pb.EncodeBatchLinksRequest{}
	body := []*pb.EncodeBatchLinksRequestEntity{}
	body = append(body, &pb.EncodeBatchLinksRequestEntity{CorrelationID: "djfjd", OriginalURL: "https://yandex.com"})
	body = append(body, &pb.EncodeBatchLinksRequestEntity{CorrelationID: "sdfds", OriginalURL: "https://google.com"})
	body = append(body, &pb.EncodeBatchLinksRequestEntity{CorrelationID: "eiuryti", OriginalURL: "https://lenta.com"})
	req.Links = body

	resp, err := s.BatchEncodeLink(ctx, &req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Link)

	links := resp.GetLink()
	assert.Len(t, links, 6)

	resp, err = s.BatchEncodeLink(ctx, &req)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, status.Error(codes.Internal, "failed batch encode link"))

}

// TestServer_UserLinks тест получения всех ссылок пользователя
func TestServer_UserLinks(t *testing.T) {
	s := Server{}
	err := config.Initialize()
	assert.NoError(t, err)
	conf := config.GetParams()
	err = repo.Initialize(conf)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ok, err := repo.CreateUser(ctx, "login", "secret")
	assert.NoError(t, err)
	assert.True(t, ok)

	userID, err := repo.AuthUser(ctx, "login", "secret")
	assert.NoError(t, err)
	token, err := jwttoken.GenerateToken(userID, conf.GetJWTSecret())
	assert.NoError(t, err)

	links := make(map[string]string)
	links["ywtyee"] = "https://ywtyee.ru"
	links["ewbmna"] = "https://ewbmna.ru"
	links["mnbxzc"] = "https://mnbxzc.ru"
	links["nmzbxc"] = "https://nmzbxc.ru"
	err = repo.PutBatchLinksArray(ctx, links, int(userID))
	assert.NoError(t, err)
	cancel()

	header := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewIncomingContext(context.Background(), header)

	req := pb.UserLinksRequest{}
	resp, err := s.UserLinks(ctx, &req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	respLinks := resp.GetLink()
	log.Println(respLinks)
	assert.Len(t, respLinks, 4)
}

// TestServer_DeleteLink тест удаления ссылок пользователя
func TestServer_DeleteLink(t *testing.T) {
	s := Server{}
	err := config.Initialize()
	assert.NoError(t, err)
	conf := config.GetParams()
	err = repo.Initialize(conf)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ok, err := repo.CreateUser(ctx, "login", "secret")
	assert.NoError(t, err)
	assert.True(t, ok)

	userID, err := repo.AuthUser(ctx, "login", "secret")
	assert.NoError(t, err)
	token, err := jwttoken.GenerateToken(userID, conf.GetJWTSecret())
	assert.NoError(t, err)

	links := make(map[string]string)
	links["ywtyee"] = "https://ywtyee.ru"
	links["ewbmna"] = "https://ewbmna.ru"
	links["mnbxzc"] = "https://mnbxzc.ru"
	links["nmzbxc"] = "https://nmzbxc.ru"
	err = repo.PutBatchLinksArray(ctx, links, int(userID))
	assert.NoError(t, err)
	cancel()

	header := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewIncomingContext(context.Background(), header)

	repData := []string{"ywtyee", "ewbmna"}
	req := pb.DeleteUserLinksRequest{Short: repData}
	resp, err := s.DeleteUserLinks(ctx, &req)
	assert.NoError(t, err)
	assert.Equal(t, uint32(202), resp.ResultCode)

	userLinks, err := repo.GetUserLinks(ctx, int(userID))
	assert.NoError(t, err)
	log.Println(userLinks)

}

// TestServer_InternalStats тест получения статистики сервера
func TestServer_InternalStats(t *testing.T) {
	s := Server{}
	err := config.Initialize()
	assert.NoError(t, err)
	conf := config.GetParams()
	err = repo.Initialize(conf)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ok, err := repo.CreateUser(ctx, "login", "secret")
	assert.NoError(t, err)
	assert.True(t, ok)

	userID, err := repo.AuthUser(ctx, "login", "secret")
	assert.NoError(t, err)

	links := make(map[string]string)
	links["ywtyee"] = "https://ywtyee.ru"
	links["ewbmna"] = "https://ewbmna.ru"
	links["mnbxzc"] = "https://mnbxzc.ru"
	links["nmzbxc"] = "https://nmzbxc.ru"
	err = repo.PutBatchLinksArray(ctx, links, int(userID))
	assert.NoError(t, err)
	cancel()

	ctx = metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{"x-real-ip": conf.GetTrustedSubnet()}))
	req := pb.InternalStatsRequest{}
	resp, err := s.InternalStats(ctx, &req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint64(1), resp.Users)
	assert.Equal(t, uint64(4), resp.Urls)
}
