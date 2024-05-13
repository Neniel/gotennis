package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Neniel/gotennis/database"
	"github.com/Neniel/gotennis/entity"
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
			uc := &getCategoryUsecase{
				DBReader: tt.fields.DBReader,
			}
			got, err := uc.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCategoryUsecase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCategoryUsecase.Get() = %v, want %v", got, tt.want)
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
			uc := &getCategoryUsecase{
				DBReader: tt.fields.DBReader,
			}
			got, err := uc.Get(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCategoryUsecase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCategoryUsecase.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
