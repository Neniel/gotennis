package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/Neniel/gotennis/database"
	"github.com/Neniel/gotennis/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/mock/gomock"
)

func TestUpdatePlayerRequest_Validate(t *testing.T) {
	id := primitive.NewObjectID()
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
		{
			name: "Request_for_updating_player_is_valid",
			fields: fields{
				ID:           id.Hex(),
				GovernmentID: "AR-1234567890",
				FirstName:    "Bob",
				MiddleName:   "Sponge",
				LastName:     "Square Pants",
				Category:     nil,
				Birthdate:    nil,
				PhoneNumber:  "+00 000 000 000",
				Email:        "bobsponge@test.com",
				Alias:        "bob",
			},
			args: args{
				id: id.Hex(),
			},
			wantErr: false,
		},
		{
			name: "Request_for_updating_player_is_not_valid_due_empty_id_in_request",
			fields: fields{
				ID:           "",
				GovernmentID: "AR-1234567890",
				FirstName:    "Bob",
				MiddleName:   "Sponge",
				LastName:     "Square Pants",
				Category:     nil,
				Birthdate:    nil,
				PhoneNumber:  "+00 000 000 000",
				Email:        "bobsponge@test.com",
				Alias:        "bob",
			},
			args: args{
				id: id.Hex(),
			},
			wantErr: true,
		},
		{
			name: "Request_for_updating_player_is_not_valid_due_empty_id_in_request_url",
			fields: fields{
				ID:           id.Hex(),
				GovernmentID: "AR-1234567890",
				FirstName:    "Bob",
				MiddleName:   "Sponge",
				LastName:     "Square Pants",
				Category:     nil,
				Birthdate:    nil,
				PhoneNumber:  "+00 000 000 000",
				Email:        "bobsponge@test.com",
				Alias:        "bob",
			},
			args: args{
				id: "",
			},
			wantErr: true,
		},
		{
			name: "Request_for_updating_player_is_not_valid_due_to_different_id",
			fields: fields{
				ID:           id.Hex(),
				GovernmentID: "AR-1234567890",
				FirstName:    "Bob",
				MiddleName:   "Sponge",
				LastName:     "Square Pants",
				Category:     nil,
				Birthdate:    nil,
				PhoneNumber:  "+00 000 000 000",
				Email:        "bobsponge@test.com",
				Alias:        "bob",
			},
			args: args{
				id: primitive.NewObjectID().Hex(),
			},
			wantErr: true,
		},
		{
			name: "Request_for_updating_player_is_not_valid_due_empty_government_id",
			fields: fields{
				ID:           id.Hex(),
				GovernmentID: "",
				FirstName:    "Bob",
				MiddleName:   "Sponge",
				LastName:     "Square Pants",
				Category:     nil,
				Birthdate:    nil,
				PhoneNumber:  "+00 000 000 000",
				Email:        "bobsponge@test.com",
				Alias:        "bob",
			},
			args: args{
				id: id.Hex(),
			},
			wantErr: true,
		},
		{
			name: "Request_for_updating_player_is_not_valid_due_empty_government_id",
			fields: fields{
				ID:           id.Hex(),
				GovernmentID: "",
				FirstName:    "Bob",
				MiddleName:   "Sponge",
				LastName:     "Square Pants",
				Category:     nil,
				Birthdate:    nil,
				PhoneNumber:  "+00 000 000 000",
				Email:        "bobsponge@test.com",
				Alias:        "bob",
			},
			args: args{
				id: id.Hex(),
			},
			wantErr: true,
		},
		{
			name: "Request_for_updating_player_is_not_valid_due_empty_first_name",
			fields: fields{
				ID:           id.Hex(),
				GovernmentID: "AR-1234567890",
				FirstName:    "",
				MiddleName:   "Sponge",
				LastName:     "Square Pants",
				Category:     nil,
				Birthdate:    nil,
				PhoneNumber:  "+00 000 000 000",
				Email:        "bobsponge@test.com",
				Alias:        "bob",
			},
			args: args{
				id: id.Hex(),
			},
			wantErr: true,
		},
		{
			name: "Request_for_updating_player_is_not_valid_due_empty_last_name",
			fields: fields{
				ID:           id.Hex(),
				GovernmentID: "AR-1234567890",
				FirstName:    "Bob",
				MiddleName:   "Sponge",
				LastName:     "",
				Category:     nil,
				Birthdate:    nil,
				PhoneNumber:  "+00 000 000 000",
				Email:        "bobsponge@test.com",
				Alias:        "bob",
			},
			args: args{
				id: id.Hex(),
			},
			wantErr: true,
		},
		{
			name: "Request_for_updating_player_is_not_valid_due_empty_email",
			fields: fields{
				ID:           id.Hex(),
				GovernmentID: "AR-1234567890",
				FirstName:    "Bob",
				MiddleName:   "Sponge",
				LastName:     "Square Pants",
				Category:     nil,
				Birthdate:    nil,
				PhoneNumber:  "+00 000 000 000",
				Email:        "",
				Alias:        "bob",
			},
			args: args{
				id: id.Hex(),
			},
			wantErr: true,
		},
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

func Test_updatePlayerUsecase_UpdatePlayer_Success(t *testing.T) {
	id := primitive.NewObjectID()
	dbReader := database.NewMockDBReader(gomock.NewController(t))
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))

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
		name           string
		fields         fields
		args           args
		prepareUsecase func()
		want           *entity.Player
		wantErr        bool
	}{
		{
			name: "Update_player",
			fields: fields{
				internalUpdatePlayerUsecases: &internalUpdatePlayerUsecases{
					ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReader),
					ValidateEmail:        NewValidateEmailUsecase(dbReader),
					ValidateAlias:        NewValidateAliasUsecase(dbReader),
				},
				DBWriter: dbWriter,
				DBReader: dbReader,
			},
			args: args{
				ctx: context.Background(),
				id:  id.Hex(),
				request: &UpdatePlayerRequest{
					ID:           id.Hex(),
					GovernmentID: "AR-1234567890",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					Category:     nil,
					Birthdate:    nil,
					PhoneNumber:  "+54 000 000 000",
					Email:        "bobsponge@test.com",
					Alias:        "bob",
				},
			},
			prepareUsecase: func() {
				dbReader.EXPECT().IsAvailable(gomock.Any(), "government_id", "AR-1234567890").Return(true, nil)
				dbReader.EXPECT().IsAvailable(gomock.Any(), "email", "bobsponge@test.com").Return(true, nil)
				dbReader.EXPECT().IsAvailable(gomock.Any(), "alias", "bob").Return(true, nil)
				dbReader.EXPECT().GetPlayer(gomock.Any(), id.Hex()).Return(&entity.Player{
					ID:           id,
					GovernmentID: "AR-1234567890",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					Category:     nil,
					Birthdate:    nil,
					PhoneNumber:  "+00 000 000 000",
					Email:        "bobsponge@test.com",
					Alias:        "bob",
				}, nil)

				dbWriter.EXPECT().UpdatePlayer(gomock.Any(), &entity.Player{
					ID:           id,
					GovernmentID: "AR-1234567890",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					Category:     nil,
					Birthdate:    nil,
					PhoneNumber:  "+54 000 000 000",
					Email:        "bobsponge@test.com",
					Alias:        "bob",
				}).Return(&entity.Player{
					ID:           id,
					GovernmentID: "AR-1234567890",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					Category:     nil,
					Birthdate:    nil,
					PhoneNumber:  "+54 000 000 000",
					Email:        "bobsponge@test.com",
					Alias:        "bob",
				}, nil)
			},
			want: &entity.Player{
				ID:           id,
				GovernmentID: "AR-1234567890",
				FirstName:    "Bob",
				MiddleName:   "Sponge",
				LastName:     "Square Pants",
				Category:     nil,
				Birthdate:    nil,
				PhoneNumber:  "+54 000 000 000",
				Email:        "bobsponge@test.com",
				Alias:        "bob",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareUsecase()
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
