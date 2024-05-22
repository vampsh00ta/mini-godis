package tests

import (
	"avito_intern/internal/errs"
	"avito_intern/internal/models"
	"bytes"
	"fmt"
	mock_service "mini-godis./internal/service/mocks"
	transport "mini-godis./internal/transport/http"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

func TestBanner_GetBannerForUser(t *testing.T) {
	tests := []struct {
		name         string
		inputParams  string
		mockFuncs    []MockMethod
		expectedCode int
		expectedBody string
	}{
		{
			name:        "ok",
			inputParams: "?tag_id=1&feature_id=1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.OrdinaryUser, models.Admin},
					[]any{true, nil},
				},
				{
					"GetBannerForUser",
					[]any{mock.Anything, int32(1), int32(1), false},
					[]any{
						models.Banner{
							Content:  "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
							IsActive: true,
						},

						nil,
					},
				},
			},
			expectedCode: 200,
			expectedBody: "{\"content\":\"{\\\"title\\\": \\\"some_title\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"}",
		},

		{
			name:        "ok with last revision",
			inputParams: "?tag_id=1&feature_id=1&use_last_revision=true",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.OrdinaryUser, models.Admin},
					[]any{true, nil},
				},
				{
					"GetBannerForUser",
					[]any{mock.Anything, int32(1), int32(1), true},
					[]any{
						models.Banner{
							Content:  "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
							IsActive: true,
						},

						nil,
					},
				},
			},
			expectedCode: 200,
			expectedBody: "{\"content\":\"{\\\"title\\\": \\\"some_title\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"}",
		},
		{
			name:        "ok ",
			inputParams: "?tag_id=1&feature_id=1&use_last_revision=true",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.OrdinaryUser, models.Admin},
					[]any{true, nil},
				},
				{
					"GetBannerForUser",
					[]any{mock.Anything, int32(1), int32(1), true},
					[]any{
						models.Banner{
							Content:  "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
							IsActive: false,
						},
						nil,
					},
				},
			},
			expectedCode: 404,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.NoRowsInResult),
		},
		{
			name:        "nil feature_id",
			inputParams: "?tag_id=1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.OrdinaryUser, models.Admin},
					[]any{true, nil},
				},
			},
			expectedCode: 400,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Validation),
		},
		{
			name:        "nil tag_id",
			inputParams: "?feature_id=1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.OrdinaryUser, models.Admin},
					[]any{true, nil},
				},
			},
			expectedCode: 400,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Validation),
		},
		{
			name:        "not logged",
			inputParams: "?tag_id=1&feature_id=1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.OrdinaryUser, models.Admin},
					[]any{false, errs.NotLogged},
				},
			},
			expectedCode: 401,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.NotLogged),
		},
		{
			name:        "auth error",
			inputParams: "?tag_id=1&feature_id=1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.OrdinaryUser, models.Admin},
					[]any{false, errs.Auth},
				},
			},
			expectedCode: 401,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Auth),
		},
		{
			name:        "auth error",
			inputParams: "?tag_id=1&feature_id=1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.OrdinaryUser, models.Admin},
					[]any{true, nil},
				},
				{
					"GetBannerForUser",
					[]any{mock.Anything, int32(1), int32(1), false},
					[]any{
						models.Banner{},
						errs.Unknown,
					},
				},
			},
			expectedCode: 500,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Unknown),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srvc := mock_service.NewService(t)
			for _, mockFunc := range test.mockFuncs {
				srvc.On(mockFunc.methodName, mockFunc.args...).Once().Return(mockFunc.returns...)
			}

			r := transport.NewTest(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			router := http.NewServeMux()
			router.HandleFunc("GET /user_banner", r.GetBannerForUser)
			req := httptest.NewRequest("GET", "/user_banner"+test.inputParams, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestBanner_GetBanners(t *testing.T) {
	tests := []struct {
		name         string
		inputParams  string
		mockFuncs    []MockMethod
		expectedCode int
		expectedBody string
	}{
		{
			name:        "ok by tag_id",
			inputParams: "?tag_id=1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"GetBanners",
					[]any{mock.Anything, int32(1), int32(0), int32(0), int32(0)},
					[]any{
						[]models.Banner{
							{
								ID:        1,
								Content:   "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
								IsActive:  true,
								UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
								CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
								Feature:   1,
								Tags: []int32{
									1, 2, 3,
								},
							},
						},

						nil,
					},
				},
			},
			expectedCode: 200,
			expectedBody: `[{"banner_id":1,"tag_ids":[1,2,3],"feature_id":1,"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}","is_active":true,"created_at":"2000-01-01T00:00:00Z","updated_at":"2000-01-01T00:00:00Z"}]`,
			// expectedBody: "{\"content\":\"{\\\"title\\\": \\\"some_title\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"}",
		},
		{
			name:        "ok by feature_id",
			inputParams: "?feature_id=1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"GetBanners",
					[]any{mock.Anything, int32(0), int32(1), int32(0), int32(0)},
					[]any{
						[]models.Banner{
							{
								ID:        1,
								Content:   "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
								IsActive:  true,
								UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
								CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
								Feature:   1,
								Tags: []int32{
									1, 2, 3,
								},
							},
						},

						nil,
					},
				},
			},
			expectedCode: 200,
			expectedBody: `[{"banner_id":1,"tag_ids":[1,2,3],"feature_id":1,"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}","is_active":true,"created_at":"2000-01-01T00:00:00Z","updated_at":"2000-01-01T00:00:00Z"}]`,
		},
		{
			name:        "not logged",
			inputParams: "",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, errs.NotLogged},
				},
			},
			expectedCode: 401,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.NotLogged),
			// expectedBody: "{\"content\":\"{\\\"title\\\": \\\"some_title\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"}",
		},
		{
			name:        "auth error",
			inputParams: "",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, errs.Auth},
				},
			},
			expectedCode: 401,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Auth),
		},
		{
			name:        "wrong  role",
			inputParams: "",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{false, nil},
				},
			},
			expectedCode: 403,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.WrongRole),
			// expectedBody: "{\"content\":\"{\\\"title\\\": \\\"some_title\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"}",
		},
		{
			name:        "nil",
			inputParams: "?feature_id=1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"GetBanners",
					[]any{mock.Anything, int32(0), int32(1), int32(0), int32(0)},
					[]any{
						[]models.Banner{{}},

						pgx.ErrNoRows,
					},
				},
			},
			expectedCode: 404,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.NoRowsInResult),
		},
		{
			name:        "server error ",
			inputParams: "",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"GetBanners",
					[]any{mock.Anything, int32(0), int32(0), int32(0), int32(0)},
					[]any{
						[]models.Banner{},

						errs.Unknown,
					},
				},
			},
			expectedCode: 500,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Unknown),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srvc := mock_service.NewService(t)
			for _, mockFunc := range test.mockFuncs {
				srvc.On(mockFunc.methodName, mockFunc.args...).Once().Return(mockFunc.returns...)
			}

			r := transport.NewTest(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			router := http.NewServeMux()
			router.HandleFunc("GET /banner", r.GetBanners)
			req := httptest.NewRequest("GET", "/banner"+test.inputParams, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestBanner_DeleteBannerByID(t *testing.T) {
	tests := []struct {
		name         string
		inputParams  string
		mockFuncs    []MockMethod
		expectedCode int
		expectedBody string
	}{
		{
			name:        "ok",
			inputParams: "/1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"DeleteBannerByID",
					[]any{mock.Anything, 1},
					[]any{
						nil,
					},
				},
			},
			expectedCode: 204,
			expectedBody: ``,
		},

		{
			name:        "not logged",
			inputParams: "/1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, errs.NotLogged},
				},
			},
			expectedCode: 401,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.NotLogged),
			// expectedBody: "{\"content\":\"{\\\"title\\\": \\\"some_title\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"}",
		},
		{
			name:        "auth error",
			inputParams: "/1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, errs.Auth},
				},
			},
			expectedCode: 401,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Auth),
			// expectedBody: "{\"content\":\"{\\\"title\\\": \\\"some_title\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"}",
		},
		{
			name:        "wrong  role",
			inputParams: "/1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{false, nil},
				},
			},
			expectedCode: 403,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.WrongRole),
			// expectedBody: "{\"content\":\"{\\\"title\\\": \\\"some_title\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"}",
		},
		{
			name:        "nil ",
			inputParams: "/1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"DeleteBannerByID",
					[]any{mock.Anything, 1},
					[]any{
						pgx.ErrNoRows,
					},
				},
			},
			expectedCode: 404,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.NoRowsInResult),
		},
		{
			name:        "server error",
			inputParams: "/1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"DeleteBannerByID",
					[]any{mock.Anything, 1},
					[]any{
						errs.Unknown,
					},
				},
			},
			expectedCode: 500,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Unknown),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srvc := mock_service.NewService(t)
			for _, mockFunc := range test.mockFuncs {
				srvc.On(mockFunc.methodName, mockFunc.args...).Once().Return(mockFunc.returns...)
			}

			r := transport.NewTest(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /banner/{id}", r.DeleteBannerByID)
			req := httptest.NewRequest("DELETE", "/banner"+test.inputParams, nil)
			mux.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestBanner_CreateBanner(t *testing.T) {
	tests := []struct {
		name         string
		inputBody    string
		mockFuncs    []MockMethod
		expectedCode int
		expectedBody string
	}{
		{
			name: "ok",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"CreateBanner",
					[]any{mock.Anything, models.Banner{
						Tags:     []int32{1},
						Feature:  1,
						Content:  "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						IsActive: true,
					}},
					[]any{
						1,
						nil,
					},
				},
			},
			expectedCode: 201,
			expectedBody: `{"id":1}`,
		},
		{
			name: "validation error",
			inputBody: `{
						"tag_ids":[1],
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
			},
			expectedCode: 400,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Validation),
		},
		{
			name: "validation error json",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"asdasd{'asdfasdf'}",
						"is_active":true
						}`,

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
			},
			expectedCode: 400,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.IncorrectJSON),
		},
		{
			name: "nil tag id error ",
			inputBody: `{
						"tag_ids":[],
						"feature_id":1,
						"content":"{}",
						"is_active":true
						}`,

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
			},
			expectedCode: 400,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Validation),
		},
		{
			name: "not logged",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{false, errs.NotLogged},
				},
			},
			expectedCode: 401,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.NotLogged),
			// expectedBody: "{\"content\":\"{\\\"title\\\": \\\"some_title\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"}",
		},
		{
			name: "auth error",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, errs.Auth},
				},
			},
			expectedCode: 401,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Auth),
		},
		{
			name: "wrong  role",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{false, nil},
				},
			},
			expectedCode: 403,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.WrongRole),
		},

		{
			name: "server error",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"CreateBanner",
					[]any{mock.Anything, models.Banner{
						Tags:     []int32{1},
						Feature:  1,
						Content:  "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						IsActive: true,
					}},
					[]any{
						-1,
						errs.Unknown,
					},
				},
			},
			expectedCode: 500,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Unknown),
		},
		{
			name: "no such tag/feature",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"CreateBanner",
					[]any{mock.Anything, models.Banner{
						Tags:     []int32{1},
						Feature:  1,
						Content:  "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						IsActive: true,
					}},
					[]any{
						-1,
						errs.NoReference,
					},
				},
			},
			expectedCode: 404,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.NoReference),
		},
		{
			name: "already exist",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"CreateBanner",
					[]any{mock.Anything, models.Banner{
						Tags:     []int32{1},
						Feature:  1,
						Content:  "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						IsActive: true,
					}},
					[]any{
						-1,
						errs.Duplicate,
					},
				},
			},
			expectedCode: 500,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Duplicate),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srvc := mock_service.NewService(t)
			for _, mockFunc := range test.mockFuncs {
				srvc.On(mockFunc.methodName, mockFunc.args...).Once().Return(mockFunc.returns...)
			}

			r := transport.NewTest(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.HandleFunc("POST /banner", r.CreateBanner)
			req := httptest.NewRequest("POST", "/banner", bytes.NewBufferString(test.inputBody))
			mux.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func initBannerChange(banner models.Banner) models.BannerChange {
	res := models.BannerChange{
		Tags:     &banner.Tags,
		Feature:  &banner.Feature,
		Content:  &banner.Content,
		IsActive: &banner.IsActive,
	}
	return res
}

func TestBanner_ChangeBanner(t *testing.T) {
	tests := []struct {
		name       string
		inputBody  string
		inputParam string

		mockFuncs    []MockMethod
		expectedCode int
		expectedBody string
	}{
		{
			name: "ok",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			inputParam: "/1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"ChangeBanner",
					[]any{mock.Anything, 1, initBannerChange(models.Banner{
						Tags:     []int32{1},
						Feature:  1,
						Content:  "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						IsActive: true,
					})},
					[]any{
						nil,
					},
				},
			},
			expectedCode: 201,
			expectedBody: ``,
		},
		//{
		//	name:       "validation error",
		//	inputParam: "/",
		//
		//	inputBody: `{
		//				"tag_ids":1,
		//				"feature_id":1,
		//				"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
		//				"is_active":true
		//				}`,
		//	mockFuncs: []MockMethod{
		//		{
		//			"Permission",
		//			[]any{mock.Anything, mock.Anything, models.Admin},
		//			[]any{true, nil},
		//		},
		//	},
		//	expectedCode: 400,
		//	expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Validation),
		//},
		{
			name: "validation error id",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			inputParam: "/asd",

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
			},
			expectedCode: 400,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.WrongID),
		},
		{
			name: "validation error json",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"asdasd{'asdfasdf'}",
						"is_active":true
						}`,
			inputParam: "/1",

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
			},
			expectedCode: 400,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.IncorrectJSON),
		},

		{
			inputParam: "/1",
			name:       "not logged",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{false, errs.NotLogged},
				},
			},
			expectedCode: 401,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.NotLogged),
			// expectedBody: "{\"content\":\"{\\\"title\\\": \\\"some_title\\\", \\\"text\\\": \\\"some_text\\\", \\\"url\\\": \\\"some_url\\\"}\"}",
		},
		{
			name:       "auth error",
			inputParam: "/1",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{false, errs.Auth},
				},
			},
			expectedCode: 401,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Auth),
		},
		{
			name: "wrong  role",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			inputParam: "/1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{false, nil},
				},
			},
			expectedCode: 403,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.WrongRole),
		},

		{
			name: "server error",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			inputParam: "/1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"ChangeBanner",
					[]any{mock.Anything, 1, initBannerChange(
						models.Banner{
							Tags:     []int32{1},
							Feature:  1,
							Content:  "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
							IsActive: true,
						})},
					[]any{
						errs.Unknown,
					},
				},
			},
			expectedCode: 500,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Unknown),
		},
		{
			name:       "no such tag/feature",
			inputParam: "/1",

			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"ChangeBanner",
					[]any{mock.Anything, 1, initBannerChange(
						models.Banner{
							Tags:     []int32{1},
							Feature:  1,
							Content:  "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
							IsActive: true,
						})},
					[]any{
						errs.NoReference,
					},
				},
			},
			expectedCode: 404,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.NoReference),
		},
		{
			name: "already exist",
			inputBody: `{
						"tag_ids":[1],
						"feature_id":1,
						"content":"{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						"is_active":true
						}`,
			inputParam: "/1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"ChangeBanner",
					[]any{mock.Anything, 1, initBannerChange(models.Banner{
						Tags:     []int32{1},
						Feature:  1,
						Content:  "{\"title\": \"some_title\", \"text\": \"some_text\", \"url\": \"some_url\"}",
						IsActive: true,
					})},
					[]any{
						errs.Duplicate,
					},
				},
			},
			expectedCode: 500,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Duplicate),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srvc := mock_service.NewService(t)
			for _, mockFunc := range test.mockFuncs {
				srvc.On(mockFunc.methodName, mockFunc.args...).Once().Return(mockFunc.returns...)
			}

			r := transport.NewTest(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.HandleFunc("PATCH /banner/{id}", r.ChangeBanner)
			req := httptest.NewRequest("PATCH", "/banner"+test.inputParam, bytes.NewBufferString(test.inputBody))
			mux.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestBanner_DeleteBannerByTagAndFeature(t *testing.T) {
	tests := []struct {
		name       string
		inputBody  string
		inputParam string

		mockFuncs    []MockMethod
		expectedCode int
		expectedBody string
	}{
		{
			name:       "ok",
			inputParam: "?tag_id=1&feature_id=1",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"DeleteBannerByTagAndFeature",
					[]any{mock.Anything, int32(1), int32(1)},
					[]any{
						1,
						nil,
					},
				},
			},
			expectedCode: 204,
			expectedBody: `{"id":1}`,
		},
		{
			name: "validation error",

			inputParam: "?tag_i=1&featur_id=1",

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
			},
			expectedCode: 400,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Validation),
		},

		{
			name:       "not logged",
			inputParam: "?tag_id=1&feature_id=1",

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{false, errs.NotLogged},
				},
			},
			expectedCode: 401,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.NotLogged),
		},
		{
			name: "auth error",
			inputBody: `{
						"tag_id":1,
						"feature_id":1
						}`,
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{false, errs.Auth},
				},
			},
			expectedCode: 401,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Auth),
		},
		{
			name:       "wrong  role",
			inputParam: "?tag_id=1&feature_id=1",

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{false, nil},
				},
			},
			expectedCode: 403,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.WrongRole),
		},

		{
			name:       "server error",
			inputParam: "?tag_id=1&feature_id=1",

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"DeleteBannerByTagAndFeature",
					[]any{mock.Anything, int32(1), int32(1)},
					[]any{
						-1,
						errs.Unknown,
					},
				},
			},
			expectedCode: 500,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Unknown),
		},
		{
			name: "no such tag/feature",

			inputParam: "?tag_id=1&feature_id=1",

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"DeleteBannerByTagAndFeature",
					[]any{mock.Anything, int32(1), int32(1)},
					[]any{
						-1,
						errs.NoRowsInResult,
					},
				},
			},
			expectedCode: 404,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.NoRowsInResult),
		},
		{
			name:       "unknown error",
			inputParam: "?tag_id=1&feature_id=1",

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"DeleteBannerByTagAndFeature",
					[]any{mock.Anything, int32(1), int32(1)},
					[]any{
						-1,
						errs.Unknown,
					},
				},
			},
			expectedCode: 500,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Unknown),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srvc := mock_service.NewService(t)
			for _, mockFunc := range test.mockFuncs {
				srvc.On(mockFunc.methodName, mockFunc.args...).Once().Return(mockFunc.returns...)
			}

			r := transport.NewTest(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /banner", r.DeleteBannerByTagAndFeature)
			req := httptest.NewRequest("DELETE", "/banner"+test.inputParam, nil)
			mux.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestBanner_GetBannerWithHistory(t *testing.T) {
	tests := []struct {
		name         string
		inputParam   string
		mockFuncs    []MockMethod
		expectedCode int
		expectedBody string
	}{
		{
			name: "ok",

			inputParam: "/1?limit=3",

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"GetBannerWithHistory",
					[]any{mock.Anything, 1, 3},
					[]any{
						[]models.Banner{
							{
								ID:        1,
								Feature:   1,
								Tags:      []int32{1},
								Content:   `"test":"1"`,
								IsActive:  true,
								UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
								CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							},
						},
						nil,
					},
				},
			},
			expectedCode: 200,
			expectedBody: `[{"banner_id":1,"tag_ids":[1],"feature_id":1,"content":"\"test\":\"1\"","is_active":true,"created_at":"2000-01-01T00:00:00Z","updated_at":"2000-01-01T00:00:00Z"}]`,
		},
		{
			name:       "validation error",
			inputParam: "/asd?limit=3",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
			},
			expectedCode: 400,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.WrongID),
		},
		{
			name:       "validation error id",
			inputParam: "/asd?limit=3",

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
			},
			expectedCode: 400,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.WrongID),
		},

		{
			inputParam: "/1?limit=3",

			name: "not logged",
			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{false, errs.NotLogged},
				},
			},
			expectedCode: 401,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.NotLogged),
		},
		{
			name:       "auth error",
			inputParam: "/1?limit=3",

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{false, errs.Auth},
				},
			},
			expectedCode: 401,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Auth),
		},
		{
			name:       "wrong  role",
			inputParam: "/1?limit=3",

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{false, nil},
				},
			},
			expectedCode: 403,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.WrongRole),
		},

		{
			name:       "server error",
			inputParam: "/1?limit=3",

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"GetBannerWithHistory",
					[]any{mock.Anything, 1, 3},
					[]any{
						[]models.Banner{
							{
								ID:        1,
								Feature:   1,
								Tags:      []int32{1},
								Content:   `"test":"1"`,
								IsActive:  true,
								UpdatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
								CreatedAt: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
							},
						},
						errs.Unknown,
					},
				},
			},
			expectedCode: 500,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.Unknown),
		},
		{
			name:       "no history",
			inputParam: "/1?limit=3",

			mockFuncs: []MockMethod{
				{
					"Permission",
					[]any{mock.Anything, mock.Anything, models.Admin},
					[]any{true, nil},
				},
				{
					"GetBannerWithHistory",
					[]any{mock.Anything, 1, 3},
					[]any{
						[]models.Banner{},
						errs.NoRowsInResult,
					},
				},
			},
			expectedCode: 404,
			expectedBody: fmt.Sprintf(`{"error":"%s"}`, errs.NoRowsInResult),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			srvc := mock_service.NewService(t)
			for _, mockFunc := range test.mockFuncs {
				srvc.On(mockFunc.methodName, mockFunc.args...).Once().Return(mockFunc.returns...)
			}

			r := transport.NewTest(srvc, LoadLoggerDev())
			w := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.HandleFunc("GET /banner_history/{id}", r.GetBannerWithHistory)
			req := httptest.NewRequest("GET", "/banner_history"+test.inputParam, nil)
			mux.ServeHTTP(w, req)
			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
