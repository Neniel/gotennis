package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Neniel/gotennis/database"
	"github.com/Neniel/gotennis/entity"
	"go.uber.org/mock/gomock"
)

func Test_createCategoryUsecase_CreateCategory_Success(t *testing.T) {
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))
	dbWriter.
		EXPECT().
		AddCategory(gomock.Any(), gomock.Any()).
		Return(
			&entity.Category{
				Name: "Category 1",
			},
			nil,
		)

	type fields struct {
		DBWriter database.DBWriter
	}
	type args struct {
		ctx     context.Context
		request *CreateCategoryRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Category
		wantErr bool
	}{
		{
			name: "Create_category",
			fields: fields{
				DBWriter: dbWriter,
			},
			args: args{
				ctx: context.TODO(),
				request: &CreateCategoryRequest{
					Name: "Category 1",
				},
			},
			want: &entity.Category{
				Name: "Category 1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &createCategoryUsecase{
				DBWriter: tt.fields.DBWriter,
			}
			got, err := uc.CreateCategory(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("createCategoryUsecase.CreateCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createCategoryUsecase.CreateCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createCategoryUsecase_CreateCategory_Failure(t *testing.T) {
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))
	dbWriter.EXPECT().AddCategory(gomock.Any(), gomock.Any()).Return(nil, errors.New("error when saving category in the db"))

	type fields struct {
		DBWriter database.DBWriter
	}
	type args struct {
		ctx     context.Context
		request *CreateCategoryRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Category
		wantErr bool
	}{
		{
			name: "Create_category_fails",
			fields: fields{
				DBWriter: dbWriter,
			},
			args: args{
				ctx: context.Background(),
				request: &CreateCategoryRequest{
					Name: "Category 1",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Create_category_with_empty_name_should_fail",
			fields: fields{
				DBWriter: dbWriter,
			},
			args: args{
				ctx: context.Background(),
				request: &CreateCategoryRequest{
					Name: "",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &createCategoryUsecase{
				DBWriter: tt.fields.DBWriter,
			}
			got, err := uc.CreateCategory(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("createCategoryUsecase.CreateCategory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createCategoryUsecase.CreateCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}
