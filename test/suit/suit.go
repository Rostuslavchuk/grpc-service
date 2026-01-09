package suit

import (
	"context"
	"net"
	"strconv"
	"testing"

	"sso/internal/config"

	sso1 "github.com/Rostuslavchuk/sso-protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcHost = "localhost"
)

type Suite struct {
	*testing.T
	Cfg        *config.Config
	AuthClient sso1.AuthClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	config := config.MustLoadByPath("../config/local.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), config.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(
		context.Background(),
		grpcAddress(config),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc server connection faild %w", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:        config,
		AuthClient: sso1.NewAuthClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}
