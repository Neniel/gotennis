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

func TestNewListCategories(t *testing.T) {
	dbReader := database.NewMockDBReader(gomock.NewController(t))
	type args struct {
		dbReader database.DBReader
	}
	tests := []struct {
		name string
		args args
		want ListCategories
	}{
		{
			name: "Create_new_list_categories_usecase",
			args: args{
				dbReader: dbReader,
			},
			want: &listCategories{
				DBReader: dbReader,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewListCategories(tt.args.dbReader); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewListCategories() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listCategories_List_Success(t *testing.T) {
	dbReader := database.NewMockDBReader(gomock.NewController(t))
	type fields struct {
		DBReader database.DBReader
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		prepareUsecase func()
		want           []entity.Category
		wantErr        bool
	}{
		{
			name: "get_a_list_of_categories",
			fields: fields{
				DBReader: dbReader,
			},
			args: args{
				ctx: context.Background(),
			},
			prepareUsecase: func() {
				dbReader.EXPECT().GetCategories(gomock.Any()).Return([]entity.Category{
					{
						Name: "Category 1",
					},
					{
						Name: "Category 2",
					},
				}, nil)
			},
			want: []entity.Category{
				{
					Name: "Category 1",
				},
				{
					Name: "Category 2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareUsecase()
			uc := &listCategories{
				DBReader: tt.fields.DBReader,
			}
			got, err := uc.Do(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("listCategories.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listCategories.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listCategories_List_Failed(t *testing.T) {
	dbReader := database.NewMockDBReader(gomock.NewController(t))
	type fields struct {
		DBReader database.DBReader
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		prepareUsecase func()
		want           []entity.Category
		wantErr        bool
	}{
		{
			name: "Error_when_getting_a_list_of_categories",
			fields: fields{
				DBReader: dbReader,
			},
			args: args{
				ctx: context.Background(),
			},
			prepareUsecase: func() {
				dbReader.EXPECT().GetCategories(gomock.Any()).Return(nil, errors.New("error when fetting the categories"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareUsecase()
			uc := &listCategories{
				DBReader: tt.fields.DBReader,
			}
			got, err := uc.Do(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("listCategories.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listCategories.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}
