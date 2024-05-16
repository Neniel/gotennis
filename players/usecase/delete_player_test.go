package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Neniel/gotennis/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestNewDeletePlayerUsecase(t *testing.T) {
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))
	type args struct {
		dbWriter database.DBWriter
	}
	tests := []struct {
		name string
		args args
		want DeletePlayerUsecase
	}{
		{
			name: "Create_new_usecase_delete_player",
			args: args{
				dbWriter: dbWriter,
			},
			want: &deletePlayerUsecase{
				DBWriter: dbWriter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeletePlayerUsecase(tt.args.dbWriter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeletePlayerUsecase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deletePlayerUsecase_Delete(t *testing.T) {
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))
	id := primitive.NewObjectID()

	type fields struct {
		DBWriter database.DBWriter
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		prepareUsecase func()
		wantErr        bool
	}{
		{
			name: "Delete_successfully",
			fields: fields{
				DBWriter: dbWriter,
			},
			args: args{
				ctx: context.Background(),
				id:  id.Hex(),
			},
			prepareUsecase: func() {
				dbWriter.EXPECT().DeletePlayer(gomock.Any(), id.Hex()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "Delete_fails",
			fields: fields{
				DBWriter: dbWriter,
			},
			args: args{
				ctx: context.Background(),
				id:  id.Hex(),
			},
			prepareUsecase: func() {
				dbWriter.EXPECT().DeletePlayer(gomock.Any(), id.Hex()).Return(errors.New("error when deleting user"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareUsecase()
			uc := &deletePlayerUsecase{
				DBWriter: tt.fields.DBWriter,
			}
			if err := uc.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("deletePlayerUsecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
