package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
	"go.uber.org/mock/gomock"
)

func TestNewUpdateCategoryUsecase(t *testing.T) {
	dbReader := database.NewMockDBReader(gomock.NewController(t))
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))
	type args struct {
		dbReader database.DBReader
		dbWriter database.DBWriter
	}
	tests := []struct {
		name string
		args args
		want UpdateCategoryUsecase
	}{
		{
			name: "Create_a_new_update_category_usecase",
			args: args{
				dbReader: dbReader,
				dbWriter: dbWriter,
			},
			want: &updateCategoryUsecase{
				DBReader: dbReader,
				DBWriter: dbWriter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUpdateCategoryUsecase(tt.args.dbReader, tt.args.dbWriter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUpdateCategoryUsecase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateCategoryRequest_Validate(t *testing.T) {
	type fields struct {
		ID   string
		Name string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ID_is_required_for_update",
			fields: fields{
				ID: "",
			},
			args: args{
				id: "",
			},
			wantErr: true,
		},
		{
			name: "ID_on_path_must_match_id_on_request_for_update",
			fields: fields{
				ID: "663d70d88264adea5d7d29bb",
			},
			args: args{
				id: "663d70d88264adea5d7d29ba",
			},
			wantErr: true,
		},
		{
			name: "Name_is_required_for_update",
			fields: fields{
				ID:   "663d70d88264adea5d7d29bb",
				Name: "",
			},
			args: args{
				id: "663d70d88264adea5d7d29bb",
			},
			wantErr: true,
		},
		{
			name: "Request_is_correct_for_update",
			fields: fields{
				ID:   "663d70d88264adea5d7d29bb",
				Name: "Category 1",
			},
			args: args{
				id: "663d70d88264adea5d7d29bb",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UpdateCategoryRequest{
				ID:   tt.fields.ID,
				Name: tt.fields.Name,
			}
			if err := r.Validate(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("UpdateCategoryRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_updateCategoryUsecase_UpdateCategory(t *testing.T) {
	dbReader := database.NewMockDBReader((gomock.NewController(t)))
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))

	type fields struct {
		DBReader database.DBReader
		DBWriter database.DBWriter
	}
	type args struct {
		ctx     context.Context
		id      string
		request *UpdateCategoryRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		prepareMocks func()
		want         *entity.Category
		wantErr      bool
	}{
		{
			name: "Validate_fails_on_empty_ID",
			fields: fields{
				DBReader: dbReader,
				DBWriter: dbWriter,
			},
			args: args{
				ctx: context.Background(),
				id:  "663d70d88264adea5d7d29bb",
				request: &UpdateCategoryRequest{
					ID:   "",
					Name: "",
				},
			},
			prepareMocks: func() {},
			want:         nil,
			wantErr:      true,
		},
		{
			name: "Validate_fails_on_different_IDs",
			fields: fields{
				DBReader: dbReader,
				DBWriter: dbWriter,
			},
			args: args{
				ctx: context.Background(),
				id:  "663d70d88264adea5d7d29bb",
				request: &UpdateCategoryRequest{
					ID:   "663d70d88264adea5d7d29ba",
					Name: "",
				},
			},
			prepareMocks: func() {},
			want:         nil,
			wantErr:      true,
		},
		{
			name: "Validate_fails_on_empty_Name",
			fields: fields{
				DBReader: dbReader,
				DBWriter: dbWriter,
			},
			args: args{
				ctx: context.Background(),
				id:  "663d70d88264adea5d7d29bb",
				request: &UpdateCategoryRequest{
					ID:   "663d70d88264adea5d7d29bb",
					Name: "",
				},
			},
			prepareMocks: func() {},
			want:         nil,
			wantErr:      true,
		},
		{
			name: "Fails_when_fetching_category",
			fields: fields{
				DBReader: dbReader,
				DBWriter: dbWriter,
			},
			args: args{
				ctx: context.Background(),
				id:  "663d70d88264adea5d7d29bb",
				request: &UpdateCategoryRequest{
					ID:   "663d70d88264adea5d7d29bb",
					Name: "Category 1",
				},
			},
			prepareMocks: func() {
				dbReader.EXPECT().GetCategory(gomock.Any(), "663d70d88264adea5d7d29bb").Return(nil, errors.New("error when fetching category"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Fails_when_fetching_category",
			fields: fields{
				DBReader: dbReader,
				DBWriter: dbWriter,
			},
			args: args{
				ctx: context.Background(),
				id:  "663d70d88264adea5d7d29bb",
				request: &UpdateCategoryRequest{
					ID:   "663d70d88264adea5d7d29bb",
					Name: "Category 1",
				},
			},
			prepareMocks: func() {
				dbReader.EXPECT().GetCategory(gomock.Any(), "663d70d88264adea5d7d29bb").Return(&entity.Category{
					Name: "1 category",
				}, nil)

				dbWriter.EXPECT().UpdateCategory(gomock.Any(), gomock.Any()).Return(nil, errors.New("error when updating the category"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Fails_when_fetching_category",
			fields: fields{
				DBReader: dbReader,
				DBWriter: dbWriter,
			},
			args: args{
				ctx: context.Background(),
				id:  "663d70d88264adea5d7d29bb",
				request: &UpdateCategoryRequest{
					ID:   "663d70d88264adea5d7d29bb",
					Name: "Category 1",
				},
			},
			prepareMocks: func() {
				dbReader.EXPECT().GetCategory(gomock.Any(), "663d70d88264adea5d7d29bb").Return(&entity.Category{
					Name: "1 category",
				}, nil)

				dbWriter.EXPECT().UpdateCategory(gomock.Any(), gomock.Any()).Return(&entity.Category{
					Name: "Category 1",
				}, nil)
			},
			want: &entity.Category{
				Name: "Category 1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks()
			uc := &updateCategoryUsecase{
				DBReader: tt.fields.DBReader,
				DBWriter: tt.fields.DBWriter,
			}
			got, err := uc.UpdateCategory(tt.args.ctx, tt.args.id, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateCategoryUsecase.UpdateCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateCategoryUsecase.UpdateCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}
