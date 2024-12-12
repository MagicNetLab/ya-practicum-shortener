package server

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-shortener/internal/config"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/jwttoken"
	"github.com/MagicNetLab/ya-practicum-shortener/internal/service/logger"
)

const (
	userRegisterMethod    = "UserRegister"
	userAuthMethod        = "UserAuth"
	encodeLinkMethod      = "EncodeLink"
	batchEncodeLinkMethod = "BatchEncodeLink"
	userLinksMethod       = "UserLinks"
	deleteUserLinksMethod = "DeleteUserLinks"
	internalStatsMethod   = "InternalStats"
	dbPingMethod          = "DBPing"
	trustedSubnetHeader   = "x-real-ip"
)

var guestMethods = []string{
	userRegisterMethod,
	userAuthMethod,
	internalStatsMethod,
	dbPingMethod,
}

var authMethods = []string{
	encodeLinkMethod,
	batchEncodeLinkMethod,
	userLinksMethod,
	deleteUserLinksMethod,
}

var trustedNetworkMethods = []string{
	internalStatsMethod,
	dbPingMethod,
}

var loggerMethods = []string{
	userRegisterMethod,
	userAuthMethod,
	encodeLinkMethod,
	batchEncodeLinkMethod,
	userLinksMethod,
	deleteUserLinksMethod,
	internalStatsMethod,
	dbPingMethod,
}

// GuestInterceptor допускает к методу только не авторизированных
func GuestInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if !validateInterceptorMethod(info.FullMethod, guestMethods) {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	var token string
	appConfig := config.GetParams()

	values := md.Get("token")
	if len(values) == 1 {
		token = values[0]
	}

	if token != "" && jwttoken.ValidateToken(token, appConfig.GetJWTSecret()) {
		return nil, status.Errorf(codes.PermissionDenied, "forbidden for authorized users")
	}

	return handler(ctx, req)
}

// AuthInterceptor допускает к методу только авторизованных пользователей
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if !validateInterceptorMethod(info.FullMethod, authMethods) {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	var token string
	appConfig := config.GetParams()

	values := md.Get("token")
	if len(values) == 1 {
		token = values[0]
	}

	if token != "" && jwttoken.ValidateToken(token, appConfig.GetJWTSecret()) {
		return handler(ctx, req)
	}

	return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
}

// LoggerInterceptor логирует все запросы
func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if !validateInterceptorMethod(info.FullMethod, loggerMethods) {
		return handler(ctx, req)
	}

	start := time.Now()

	resp, err := handler(ctx, req)

	duration := time.Since(start)

	args := map[string]interface{}{
		"method":   info.FullMethod,
		"duration": duration,
	}
	logger.Info("grpc request info", args)

	return resp, err
}

// TrustedNetworkInterceptor допускает к методу только запросы из доверенных сетей
func TrustedNetworkInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if !validateInterceptorMethod(info.FullMethod, trustedNetworkMethods) {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	realIP := md.Get(trustedSubnetHeader)
	if len(realIP) == 0 {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	ip := realIP[0]
	conf := config.GetParams()
	trustedSubnet := conf.GetTrustedSubnet()
	if ip == "" || trustedSubnet == "" || ip != trustedSubnet {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	return handler(ctx, req)
}

func validateInterceptorMethod(method string, allows []string) bool {
	for _, allow := range allows {
		if strings.Contains(method, allow) {
			return true
		}
	}

	return false
}
