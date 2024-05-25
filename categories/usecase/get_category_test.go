package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func Test_getCategoryUsecase_Get_Success(t *testing.T) {
	dbreader := database.NewMockDBReader(gomock.NewController(t))
	dbreader.EXPECT().GetCategory(gomock.Any(), "663d70d88264adea5d7d29bb").Return(&entity.Category{
		ID: func() primitive.ObjectID {
			id, _ := primitive.ObjectIDFromHex("663d70d88264adea5d7d29bb")
			return id
		}(),
		Name: "1ra categoría",
	}, nil)

	type fields struct {
		DBReader database.DBReader
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Category
		wantErr bool
	}{
		{
			name: "Get_Category",
			fields: fields{
				DBReader: dbreader,
			},
			args: args{
				ctx: context.TODO(),
				id:  "663d70d88264adea5d7d29bb",
			},
			want: &entity.Category{
				ID: func() primitive.ObjectID {
					id, _ := primitive.ObjectIDFromHex("663d70d88264adea5d7d29bb")
					return id
				}(),
				Name: "1ra categoría",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &getCategory{
				DBReader: tt.fields.DBReader,
			}
			got, err := uc.Do(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCategory.Do error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCategory.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCategoryUsecase_Get_Failure(t *testing.T) {
	dbreader := database.NewMockDBReader(gomock.NewController(t))
	dbreader.EXPECT().GetCategory(gomock.Any(), "663d70d88264adea5d7d29bb").Return(nil, errors.New("error when fetching category"))

	type fields struct {
		DBReader database.DBReader
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Category
		wantErr bool
	}{
		{
			name: "Get_Category",
			fields: fields{
				DBReader: dbreader,
			},
			args: args{
				ctx: context.TODO(),
				id:  "663d70d88264adea5d7d29bb",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &getCategory{
				DBReader: tt.fields.DBReader,
			}
			got, err := uc.Do(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCategory.Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCategory.Do() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGetCategoryUsecase(t *testing.T) {
	dbReader := database.NewMockDBReader(gomock.NewController(t))
	type args struct {
		dbReader database.DBReader
	}
	tests := []struct {
		name string
		args args
		want GetCategory
	}{
		{
			name: "New_get_category_usecase",
			args: args{dbReader: dbReader},
			want: &getCategory{DBReader: dbReader},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGetCategory(tt.args.dbReader); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGetCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}
