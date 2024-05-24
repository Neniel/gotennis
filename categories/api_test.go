package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Neniel/gotennis/lib/app"
	"go.uber.org/mock/gomock"
)

func TestAPIServer_pingHandler(t *testing.T) {
	type fields struct {
		CategoryMicroservice *CategoryMicroservice
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Should return 200 OK",
			fields: fields{
				CategoryMicroservice: &CategoryMicroservice{
					App:      app.NewMockIApp(gomock.NewController(t)),
					Usecases: &Usecases{},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: &http.Request{Method: http.MethodGet},
			},
			expectedStatusCode: 200,
		},
		{
			name: "Should return 405 Method not allowed",
			fields: fields{
				CategoryMicroservice: &CategoryMicroservice{
					App:      app.NewMockIApp(gomock.NewController(t)),
					Usecases: &Usecases{},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: &http.Request{Method: http.MethodPost},
			},
			expectedStatusCode: 405,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &APIServer{
				CategoryMicroservice: tt.fields.CategoryMicroservice,
			}
			api.pingHandler(tt.args.w, tt.args.r)
			if tt.args.w.Code != tt.expectedStatusCode {
				t.Fatalf("expected %d and got %d", tt.expectedStatusCode, tt.args.w.Code)
			}
		})
	}
}
