package app

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sfu-teamproject/smartbuy/backend/apperrors"
	"github.com/sfu-teamproject/smartbuy/backend/logger/mocklogger"
	"github.com/sfu-teamproject/smartbuy/backend/models"
	"github.com/sfu-teamproject/smartbuy/backend/storage/mockstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExtractPathValue(t *testing.T) {
	testApp := NewApp(nil, nil, nil)

	tests := []struct {
		name        string
		pathValue   string
		pathKey     string
		setupReq    func(*http.Request)
		expectedID  int
		expectError bool
		errorType   error
	}{
		{
			name:      "Valid ID",
			pathValue: "123",
			pathKey:   "id",
			setupReq: func(r *http.Request) {
				r.SetPathValue("id", "123")
			},
			expectedID:  123,
			expectError: false,
		},
		{
			name:      "Empty ID string",
			pathValue: "",
			pathKey:   "id",
			setupReq: func(r *http.Request) {
				r.SetPathValue("id", "")
			},
			expectedID:  0,
			expectError: true,
			errorType:   apperrors.ErrBadRequest,
		},
		{
			name:      "Non-integer ID",
			pathValue: "abc",
			pathKey:   "id",
			setupReq: func(r *http.Request) {
				r.SetPathValue("id", "abc")
			},
			expectedID:  0,
			expectError: true,
			errorType:   apperrors.ErrBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.setupReq != nil {
				tc.setupReq(req)
			}

			id, err := testApp.ExtractPathValue(req, tc.pathKey)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorType != nil {
					assert.True(t, errors.Is(err, tc.errorType), "Expected error type %v, got %v", tc.errorType, err)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedID, id)
			}
		})
	}
}

type MockClaims struct {
	jwt.RegisteredClaims
	Role models.Role
}

func (m *MockClaims) GetSubject() (string, error) {
	if m.Subject == "" {
		return "", errors.New("subject missing")
	}
	return m.Subject, nil
}

func (m *MockClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	if m.ExpiresAt == nil {
		return nil, errors.New("expiration time missing")
	}
	return m.ExpiresAt, nil
}

func (m *MockClaims) GetIssuer() (string, error) {
	if m.Issuer == "" {
		return "", errors.New("issuer missing")
	}
	return m.Issuer, nil
}

// --- NewApp Tests ---
func TestNewApp(t *testing.T) {
	// Set JWT_SECRET environment variable for the test
	t.Setenv("JWT_SECRET", "supersecretkey")

	mockLog := new(mocklogger.MockLogger)
	mockServer := &http.Server{}
	mockDB := new(mockstorage.MockStorage)

	a := NewApp(mockLog, mockServer, mockDB)

	assert.NotNil(t, a)
	assert.Equal(t, mockLog, a.Log)
	assert.Equal(t, mockServer, a.Server)
	assert.Equal(t, mockDB, a.DB)
	assert.Equal(t, []byte("supersecretkey"), a.jwtSecret)
}

// --- ErrorJSON Tests ---
func TestErrorJSON(t *testing.T) {
	tests := []struct {
		name         string
		inputError   error
		expectedCode int
		expectedBody string
		loggedError  error // The error message that should be logged
	}{
		{
			name:         "ErrNotFound",
			inputError:   apperrors.ErrNotFound,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"error":"not found"}`,
			loggedError:  apperrors.ErrNotFound,
		},
		{
			name:         "ErrBadRequest",
			inputError:   apperrors.ErrBadRequest,
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"error":"bad request"}`,
			loggedError:  apperrors.ErrBadRequest,
		},
		{
			name:         "ErrInvalidCredentials",
			inputError:   apperrors.ErrInvalidCredentials,
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"unauthorized"}`,
			loggedError:  apperrors.ErrInvalidCredentials, // original error is logged
		},
		{
			name:         "ErrUnauthorized",
			inputError:   apperrors.ErrUnauthorized,
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"error":"unauthorized"}`,
			loggedError:  apperrors.ErrUnauthorized,
		},
		{
			name:         "ErrForbidden",
			inputError:   apperrors.ErrForbidden,
			expectedCode: http.StatusForbidden,
			expectedBody: `{"error":"forbidden"}`,
			loggedError:  apperrors.ErrForbidden,
		},
		{
			name:         "ErrAlreadyExists",
			inputError:   apperrors.ErrAlreadyExists,
			expectedCode: http.StatusConflict,
			expectedBody: `{"error":"already exists"}`,
			loggedError:  apperrors.ErrAlreadyExists,
		},
		{
			name:         "Generic Error",
			inputError:   errors.New("something went wrong"),
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"error":"internal server error"}`,
			loggedError:  errors.New("something went wrong"), // original error is logged
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockLog := new(mocklogger.MockLogger)
			testApp := NewApp(mockLog, nil, nil) // Only logger matters here

			req := httptest.NewRequest(http.MethodGet, "/test", nil) // This path is fixed for all tc in ErrorJSON
			mockLog.On("Errorln", req.Method, mock.MatchedBy(func(u *url.URL) bool { return u.Path == req.URL.Path }), tc.loggedError.Error()).Return()

			rr := httptest.NewRecorder()

			testApp.ErrorJSON(rr, req, tc.inputError)

			assert.Equal(t, tc.expectedCode, rr.Code)
			assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
			assert.Equal(t, "nosniff", rr.Header().Get("X-Content-Type-Options"))
			assert.JSONEq(t, tc.expectedBody, rr.Body.String())

			mockLog.AssertExpectations(t)
		})
	}
}

// --- Encode Tests ---
func TestEncode(t *testing.T) {
	mockLog := new(mocklogger.MockLogger)
	testApp := NewApp(mockLog, nil, nil)

	t.Run("Success encoding struct", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		type TestStruct struct {
			Name  string `json:"name"`
			Value int    `json:"value"`
		}
		data := TestStruct{Name: "Item1", Value: 100}

		testApp.Encode(rr, req, data)

		assert.Equal(t, http.StatusOK, rr.Code) // Encode doesn't set status code
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))

		expectedJSON := `{
			"name": "Item1",
			"value": 100
		}`
		// Use assert.JSONEq for robust comparison of JSON strings
		assert.JSONEq(t, expectedJSON, rr.Body.String())
		mockLog.AssertNotCalled(t, "Errorln", mock.Anything, mock.Anything, mock.Anything)
	})
}

// --- GetClaims Tests ---
func TestGetClaims(t *testing.T) {
	testApp := NewApp(nil, nil, nil) // App instance not critical for this test

	// Helper to create a context with claims
	createContextWithClaims := func(userID int, role models.Role) context.Context {
		// Using app.Claims as it's the actual type expected by GetClaims
		claims := &Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   strconv.Itoa(userID),
				Issuer:    "Smartbuy",                                    // Set a dummy issuer to avoid issues in app.Auth
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)), // Dummy expiry
			},
			Role: role,
		}
		return context.WithValue(context.Background(), ClaimsKey, claims)
	}

	t.Run("Success", func(t *testing.T) {
		ctx := createContextWithClaims(123, models.RoleUser)
		req := httptest.NewRequest(http.MethodGet, "/test", nil).WithContext(ctx)

		userID, role, err := testApp.GetClaims(req)
		assert.NoError(t, err)
		assert.Equal(t, 123, userID)
		assert.Equal(t, models.RoleUser, role)
	})

	t.Run("No Claims in Context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil) // No context with claims

		userID, role, err := testApp.GetClaims(req)
		assert.Error(t, err)
		assert.Equal(t, "error getting claims", err.Error())
		assert.Equal(t, 0, userID)
		assert.Equal(t, models.Role(""), role)
	})

	t.Run("Claims Type Mismatch", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ClaimsKey, "wrong type") // Pass wrong type
		req := httptest.NewRequest(http.MethodGet, "/test", nil).WithContext(ctx)

		userID, role, err := testApp.GetClaims(req)
		assert.Error(t, err)
		assert.Equal(t, "error getting claims", err.Error())
		assert.Equal(t, 0, userID)
		assert.Equal(t, models.Role(""), role)
	})
}
