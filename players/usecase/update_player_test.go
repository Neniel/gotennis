package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/Neniel/gotennis/database"
	"github.com/Neniel/gotennis/entity"
	"go.uber.org/mock/gomock"
)

func TestUpdatePlayerRequest_Validate(t *testing.T) {
	type fields struct {
		ID           string
		GovernmentID string
		FirstName    string
		MiddleName   string
		LastName     string
		Category     *entity.Category
		Birthdate    *time.Time
		PhoneNumber  string
		Email        string
		Alias        string
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &UpdatePlayerRequest{
				ID:           tt.fields.ID,
				GovernmentID: tt.fields.GovernmentID,
				FirstName:    tt.fields.FirstName,
				MiddleName:   tt.fields.MiddleName,
				LastName:     tt.fields.LastName,
				Category:     tt.fields.Category,
				Birthdate:    tt.fields.Birthdate,
				PhoneNumber:  tt.fields.PhoneNumber,
				Email:        tt.fields.Email,
				Alias:        tt.fields.Alias,
			}
			if err := r.Validate(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePlayerRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewUpdatePlayerUsecase(t *testing.T) {
	dbReader := database.NewMockDBReader(gomock.NewController(t))
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))
	type args struct {
		dbWriter database.DBWriter
		dbReader database.DBReader
	}
	tests := []struct {
		name string
		args args
		want UpdatePlayerUsecase
	}{
		{
			name: "Create_new_update_player_usecase",
			args: args{
				dbWriter: dbWriter,
				dbReader: dbReader,
			},
			want: &updatePlayerUsecase{
				DBWriter: dbWriter,
				DBReader: dbReader,
				internalUpdatePlayerUsecases: &internalUpdatePlayerUsecases{
					ValidateGovernmentID: &validateGovernmentIDUsecase{
						DBReader: dbReader,
					},
					ValidateEmail: &validateEmailUsecaseUsecase{
						DBReader: dbReader,
					},
					ValidateAlias: &validateAliasUsecase{
						DBReader: dbReader,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUpdatePlayerUsecase(tt.args.dbWriter, tt.args.dbReader); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUpdatePlayerUsecase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updatePlayerUsecase_UpdatePlayer(t *testing.T) {
	type fields struct {
		internalUpdatePlayerUsecases *internalUpdatePlayerUsecases
		DBWriter                     database.DBWriter
		DBReader                     database.DBReader
	}
	type args struct {
		ctx     context.Context
		id      string
		request *UpdatePlayerRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Player
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &updatePlayerUsecase{
				internalUpdatePlayerUsecases: tt.fields.internalUpdatePlayerUsecases,
				DBWriter:                     tt.fields.DBWriter,
				DBReader:                     tt.fields.DBReader,
			}
			got, err := uc.UpdatePlayer(tt.args.ctx, tt.args.id, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("updatePlayerUsecase.UpdatePlayer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updatePlayerUsecase.UpdatePlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}
