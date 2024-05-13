package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Neniel/gotennis/database"
	"go.uber.org/mock/gomock"
)

func TestNewDeleteCategoryUsecase(t *testing.T) {
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))

	type args struct {
		dbWriter database.DBWriter
	}
	tests := []struct {
		name string
		args args
		want DeleteCategoryUsecase
	}{
		{
			name: "NewDeleteCategoryUsecase",
			args: args{
				dbWriter: dbWriter,
			},
			want: &deleteCategoryUsecase{
				DBWriter: dbWriter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeleteCategoryUsecase(tt.args.dbWriter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeleteCategoryUsecase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deleteCategoryUsecase_Delete_Success(t *testing.T) {
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))

	type fields struct {
		DBWriter database.DBWriter
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		prepareMocks func()
		wantErr      bool
	}{
		{
			name: "Deletes_category_successfully",
			fields: fields{
				DBWriter: dbWriter,
			},
			args: args{
				ctx: context.Background(),
				id:  "663d70d88264adea5d7d29bb",
			},
			prepareMocks: func() {
				dbWriter.EXPECT().DeleteCategory(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks()
			uc := &deleteCategoryUsecase{
				DBWriter: tt.fields.DBWriter,
			}
			if err := uc.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("deleteCategoryUsecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_deleteCategoryUsecase_Delete_Failure(t *testing.T) {
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))

	type fields struct {
		DBWriter database.DBWriter
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		prepareMocks func()
		wantErr      bool
	}{
		{
			name: "Deletes_category_fails",
			fields: fields{
				DBWriter: dbWriter,
			},
			args: args{
				ctx: context.Background(),
				id:  "663d70d88264adea5d7d29bb",
			},
			prepareMocks: func() {
				dbWriter.EXPECT().DeleteCategory(gomock.Any(), gomock.Any()).Return(errors.New("error when deleting the category"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks()
			uc := &deleteCategoryUsecase{
				DBWriter: tt.fields.DBWriter,
			}
			if err := uc.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("deleteCategoryUsecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
