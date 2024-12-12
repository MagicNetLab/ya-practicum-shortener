package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/repo/postgres"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/app/shortgen"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/jwttoken"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/link"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/user"
	pb "github.com/MagicNetLab/ya-practicum-shortener/pkg/shortener_proto"
)

// Server сервер
type Server struct {
	pb.GrpcShortenerServer
}

// UserRegister регистрация пользователя
func (s *Server) UserRegister(ctx context.Context, r *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	login := r.GetLogin()
	secret := r.GetSecret()

	exists, err := user.HasLogin(login)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed check user login", args)

		return nil, status.Error(codes.Internal, "failed registration user")
	}

	if exists {
		return nil, status.Error(codes.AlreadyExists, "login already used")
	}

	userID, err := user.Create(login, secret)
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed create user", args)

		return nil, status.Error(codes.Internal, "failed registration user")
	}

	conf := config.GetParams()
	token, err := jwttoken.GenerateToken(userID, conf.GetJWTSecret())
	if err != nil {
		args := map[string]interface{}{"error": err.Error()}
		logger.Error("failed generate token for new user", args)

		return nil, status.Error(codes.Internal, "failed generate token")
	}

	return &pb.UserRegisterResponse{Token: token}, nil

}

// UserAuth авторизация пользователя
func (s *Server) UserAuth(ctx context.Context, r *pb.UserAuthRequest) (*pb.UserAuthResponse, error) {
	userID, err := user.Authenticate(r.GetLogin(), r.GetSecret())
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed authentication")
	}

	conf := config.GetParams()
	token, err := jwttoken.GenerateToken(userID, conf.GetJWTSecret())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed generate token")
	}

	return &pb.UserAuthResponse{Token: token}, nil
}

// EncodeLink сокращение одной ссылки
func (s *Server) EncodeLink(ctx context.Context, r *pb.EncodeLinkRequest) (*pb.EncodeLinkResponse, error) {
	userID, err := getUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed parse userID from headers")
	}

	short, responseCode := link.Shorten(r.GetOriginalUrl(), int(userID))
	if responseCode == http.StatusCreated {
		return &pb.EncodeLinkResponse{ShortLink: formingShortLink(short)}, nil
	}

	if responseCode == http.StatusConflict {
		return nil, status.Error(codes.AlreadyExists, "link already exists")
	}

	return nil, status.Error(codes.Internal, "failed generate short link")
}

// BatchEncodeLink сокращение пакета ссылок
func (s *Server) BatchEncodeLink(ctx context.Context, r *pb.EncodeBatchLinksRequest) (*pb.EncodeBatchLinksResponse, error) {
	userID, err := getUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed parse userID from headers")
	}

	links := r.GetLinks()
	resultLinks := make([]*pb.EncodeBatchLinksResponseEntity, len(links))
	storeData := make(map[string]string)
	for _, l := range links {
		short := shortgen.GetShortLink(7)
		row := pb.EncodeBatchLinksResponseEntity{
			CorrelationID: l.GetCorrelationID(),
			ShortURL:      formingShortLink(short),
		}
		storeData[short] = l.GetOriginalURL()
		resultLinks = append(resultLinks, &row)
	}

	err = repo.PutBatchLinksArray(ctx, storeData, int(userID))
	if err != nil {
		if errors.Is(err, postgres.ErrLinkUniqueConflict) {
			return nil, status.Error(codes.AlreadyExists, "conflict: one or more links are not unique")
		}

		return nil, status.Error(codes.Internal, "failed batch encode link")
	}

	return &pb.EncodeBatchLinksResponse{Link: resultLinks}, nil
}

// UserLinks получение всех ссылок пользователя
func (s *Server) UserLinks(ctx context.Context, r *pb.UserLinksRequest) (*pb.UserLinksResponse, error) {
	userID, err := getUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed parse userID from headers")
	}

	res, err := repo.GetUserLinks(ctx, int(userID))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed get user links")
	}

	c := config.GetParams()
	urlPrefix := "http"
	if c.IsEnableHTTPS() {
		urlPrefix = "https"
	}

	resp := &pb.UserLinksResponse{}
	for k, v := range res {
		row := pb.UserLinksResponseEntity{
			ShortURL:    fmt.Sprintf("%s://%s/%s", urlPrefix, c.GetShortHost(), k),
			OriginalURL: v,
		}

		resp.Link = append(resp.Link, &row)
	}

	return resp, nil
}

// DeleteUserLinks удаление ссылок пользователя
func (s *Server) DeleteUserLinks(ctx context.Context, r *pb.DeleteUserLinksRequest) (*pb.DeleteUserLinksResponse, error) {
	userID, err := getUserIDFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "failed parse userID from headers")
	}

	go link.BatchDeleteLinks(ctx, r.GetShort(), int(userID))

	return &pb.DeleteUserLinksResponse{ResultCode: 202}, nil
}

// InternalStats общая статистика сервера
func (s *Server) InternalStats(ctx context.Context, r *pb.InternalStatsRequest) (*pb.InternalStatsResponse, error) {
	urls, err := repo.GetLinksCount(ctx)
	if err != nil {
		return nil, errors.New("failed get server stats")
	}

	users, err := repo.GetUsersCount(ctx)
	if err != nil {
		return nil, errors.New("failed get server stats")
	}

	return &pb.InternalStatsResponse{Urls: uint64(urls), Users: uint64(users)}, nil
}

// DBPing проверка состояния соединения с БД
func (s *Server) DBPing(ctx context.Context, r *pb.DBPingRequest) (*pb.DBPingResponse, error) {
	isPostgresOk := postgres.Ping()
	if !isPostgresOk {
		return nil, status.Error(codes.Unavailable, "failed ping")
	}

	return &pb.DBPingResponse{ResultCode: 200}, nil
}

// getUserIDFromCtx получение userID из контекста запроса
func getUserIDFromCtx(ctx context.Context) (int64, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("failed to parse metadata")
	}

	val := meta.Get("token")
	if len(val) == 0 {
		return 0, errors.New("no token provided")
	}

	token := val[0]
	userID, err := jwttoken.GetUserIDFromToken(token)
	if err != nil {
		return 0, errors.New("forbidden: no token provided")
	}

	return userID, nil
}

// formingShortLink формирование готовой короткой ссылки
func formingShortLink(short string) string {
	conf := config.GetParams()

	prefix := "http"
	if conf.IsEnableHTTPS() {
		prefix = "https"
	}

	return fmt.Sprintf("%s://%s/%s", prefix, conf.GetShortHost(), short)
}
