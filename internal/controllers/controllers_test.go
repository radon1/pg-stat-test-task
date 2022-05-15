package controllers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/radon1/pg-stat-test-task/internal/models"
	"github.com/radon1/pg-stat-test-task/internal/repositories"
	"github.com/radon1/pg-stat-test-task/internal/system/logger"
	pgstatrepomocks "github.com/radon1/pg-stat-test-task/test/mocks/packages/pgstatrepo"
)

func Test_FindQueriesStat(t *testing.T) {
	type mocks struct {
		pgStatRepository *pgstatrepomocks.MockRepository
	}

	cases := []struct {
		name         string
		route        string
		initMocks    func(mocks mocks)
		expectedCode int
		expectedBody string
	}{
		{
			name:  "without params",
			route: "/queries-stat",
			initMocks: func(mocks mocks) {
				mocks.pgStatRepository.
					EXPECT().
					FindQueriesStat(gomock.Any(), 20, 0, nil).
					Return([]models.QueryStat{
						{
							Query:      "select 1",
							Calls:      2,
							TotalTime:  12.2,
							MeanTime:   0.2,
							Percentage: 10.2,
						},
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"query":"select 1","calls":2,"total_time":12.2,"mean_time":0.2,"percentage":10.2}]`,
		},
		{
			name:  "with params",
			route: "/queries-stat?limit=10&offset=2&filter=select,delete",
			initMocks: func(mocks mocks) {
				mocks.pgStatRepository.
					EXPECT().
					FindQueriesStat(gomock.Any(), 10, 2, []string{"select", "delete"}).
					Return([]models.QueryStat{
						{
							Query:      "select 1",
							Calls:      2,
							TotalTime:  12.2,
							MeanTime:   0.2,
							Percentage: 10.2,
						},
						{
							Query:      "delete from 1",
							Calls:      2,
							TotalTime:  12.2,
							MeanTime:   0.2,
							Percentage: 10.2,
						},
					}, nil).
					Times(1)
			},
			expectedCode: http.StatusOK,
			expectedBody: `[{"query":"select 1","calls":2,"total_time":12.2,"mean_time":0.2,"percentage":10.2},{"query":"delete from 1","calls":2,"total_time":12.2,"mean_time":0.2,"percentage":10.2}]`,
		},
		{
			name:         "invalid limit",
			route:        "/queries-stat?limit=ss",
			initMocks:    func(mocks mocks) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"invalid params"}`,
		},
		{
			name:         "invalid limit",
			route:        "/queries-stat?offset=ss",
			initMocks:    func(mocks mocks) {},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"invalid params"}`,
		},
		{
			name:  "invalid limit",
			route: "/queries-stat",
			initMocks: func(mocks mocks) {
				mocks.pgStatRepository.
					EXPECT().
					FindQueriesStat(gomock.Any(), 20, 0, gomock.Any()).
					Return(nil, errors.New("custom error")).
					Times(1)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"something is wrong"}`,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			var (
				testFiberApp = fiber.New()
				mockCtrl     = gomock.NewController(t)
			)
			defer mockCtrl.Finish()

			testLogger, err := logger.New("error")
			assert.NoError(t, err)

			m := mocks{
				pgStatRepository: pgstatrepomocks.NewMockRepository(mockCtrl),
			}
			testCase.initMocks(m)

			Register(
				testLogger,
				&repositories.Container{PGStat: m.pgStatRepository},
				testFiberApp,
			)

			req := httptest.NewRequest("GET", testCase.route, nil)
			resp, err := testFiberApp.Test(req, 1)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedCode, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedBody, string(body))
		})
	}
}
