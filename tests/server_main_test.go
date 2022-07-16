//go:build all || integration
// +build all integration

package tests

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	h "gitlab.com/g6834/team26/task/internal/adapters/http"
	"gitlab.com/g6834/team26/task/internal/adapters/postgres"
	"gitlab.com/g6834/team26/task/internal/domain/task"
	"gitlab.com/g6834/team26/task/pkg/api"
	"gitlab.com/g6834/team26/task/pkg/config"
	"gitlab.com/g6834/team26/task/pkg/logger"
	"gitlab.com/g6834/team26/task/pkg/mocks"
)

type TestcontainersSuite struct {
	suite.Suite

	srv           *h.Server
	pgContainer   testcontainers.Container
	gAuthMock     *mocks.GrpcAuthMock // TODO: заменить на законченную версию сервиса auth
	gAnalyticMock *mocks.GrpcAnalyticMock
}

const (
	dbName = "mtsteta"
	dbUser = "postgres"
	dbPass = "1111"
)

func TestTestcontainers(t *testing.T) {
	suite.Run(t, &TestcontainersSuite{})
}

func (s *TestcontainersSuite) SetupSuite() {
	l := logger.New()
	ctx := context.Background()

	c, err := config.New()
	if err != nil {
		s.Suite.T().Errorf("Error parsing env: %s", err)
	}

	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:14",
			ExposedPorts: []string{"5432"},
			Env: map[string]string{
				"POSTGRES_DB":       dbName,
				"POSTGRES_USER":     dbUser,
				"POSTGRES_PASSWORD": dbPass,
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections"),
			SkipReaper: true,
			AutoRemove: true,
			Mounts: testcontainers.ContainerMounts{
				testcontainers.ContainerMount{
					Source: testcontainers.GenericBindMountSource{
						HostPath: "D:/Git/job_projects/mts_teta_projects/task/db.sql",
					},
					Target:   "/docker-entrypoint-initdb.d/db.sql",
					ReadOnly: false},
			},
		},
		Started: true,
	})
	s.Require().NoError(err)

	time.Sleep(5 * time.Second)

	dbPort, err := dbContainer.MappedPort(ctx, "5432")
	s.Require().NoError(err)
	dbIp, err := dbContainer.Host(ctx)
	s.Require().NoError(err)

	pgconn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", dbUser, dbPass, dbIp, uint16(dbPort.Int()), dbName)
	db, err := postgres.New(ctx, pgconn)
	if err != nil {
		s.Suite.T().Errorf("db init failed: %s", err)
		s.Suite.T().FailNow()
	}

	// grpcconn := getenv.GetEnv("GRPC_URL", "localhost:4000")
	// grpc, err := grpc.New(grpcconn)
	// if err != nil {
	// 	s.Suite.T().Errorf("grpc client init failed: %s", err)
	// 	s.Suite.T().FailNow()
	// }
	gAuth := new(mocks.GrpcAuthMock)
	gAnalytic := new(mocks.GrpcAnalyticMock)

	taskS := task.New(db, gAuth, gAnalytic)

	srv, err := h.New(l, taskS, c)
	if err != nil {
		s.Suite.T().Errorf("http server creating failed: %s", err)
		s.Suite.T().FailNow()
	}

	s.srv = srv
	s.pgContainer = dbContainer
	s.gAuthMock = gAuth
	s.gAnalyticMock = gAnalytic

	go s.srv.Start()

	s.T().Log("Suite setup is done")
}

func (s *TestcontainersSuite) TearDownSuite() {
	_ = s.srv.Stop(context.Background())
	s.pgContainer.Terminate(context.Background())
	s.T().Log("Suite stop is done")
}

func (s *TestcontainersSuite) TestDBSelect() {
	ctx := context.Background()
	s.gAuthMock.On("Validate", ctx, mock.Anything, mock.Anything).Return(&api.AuthResponse{Result: true, Login: "test123", AccessToken: new(api.Token), RefreshToken: new(api.Token)}, nil)
	bodyReq := strings.NewReader("")

	req, err := http.NewRequest("GET", "http://localhost:3000/task/v1/tasks/", bodyReq)
	s.NoError(err)

	client := http.Client{}
	response, err := client.Do(req)

	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)
	response.Body.Close()
}
