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

func Test_createCategory_CreateCategory_Success(t *testing.T) {
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))

	type fields struct {
		DBWriter database.DBWriter
	}
	type args struct {
		ctx     context.Context
		request *CreateCategoryRequest
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
			prepareMocks: func() {
				dbWriter.
					EXPECT().
					AddCategory(gomock.Any(), gomock.Any()).
					Return(
						&entity.Category{
							Name: "Category 1",
						},
						nil,
					)
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
			uc := &createCategory{
				DBWriter: tt.fields.DBWriter,
			}
			got, err := uc.Do(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("createCategory.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createCategory.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createCategoryUsecase_CreateCategory_Failure(t *testing.T) {
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))

	type fields struct {
		DBWriter database.DBWriter
	}
	type args struct {
		ctx     context.Context
		request *CreateCategoryRequest
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
			prepareMocks: func() {
				dbWriter.EXPECT().AddCategory(gomock.Any(), gomock.Any()).Return(nil, errors.New("error when saving category in the db"))
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
			prepareMocks: func() {},
			want:         nil,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks()
			uc := &createCategory{
				DBWriter: tt.fields.DBWriter,
			}
			got, err := uc.Do(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("createCategory.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createCategory.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCreateCategoryUsecase(t *testing.T) {
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))
	type args struct {
		dbWriter database.DBWriter
	}
	tests := []struct {
		name string
		args args
		want CreateCategory
	}{
		{
			name: "Create_category_usecase",
			args: args{
				dbWriter: dbWriter,
			},
			want: &createCategory{
				DBWriter: dbWriter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCreateCategory(tt.args.dbWriter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCreateCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}
