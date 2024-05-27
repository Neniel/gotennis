package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/util"
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
		Alias        *string
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
				Alias:        util.ToPtr("bob"),
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
				Alias:        util.ToPtr("bob"),
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
				Alias:        util.ToPtr("bob"),
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
				Alias:        util.ToPtr("bob"),
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
				Alias:        util.ToPtr("bob"),
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
				Alias:        util.ToPtr("bob"),
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
				Alias:        util.ToPtr("bob"),
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
				Alias:        util.ToPtr("bob"),
			},
			args: args{
				id: id.Hex(),
			},
			wantErr: true,
		},
		{
			name: "Request_for_updating_player_is_not_valid_due_to_zero_date",
			fields: fields{
				ID:           id.Hex(),
				GovernmentID: "AR-1234567890",
				FirstName:    "Bob",
				MiddleName:   "Sponge",
				LastName:     "Square Pants",
				Category:     nil,
				Birthdate:    &time.Time{},
				PhoneNumber:  "+00 000 000 000",
				Email:        "bobsponge@gmail.com",
				Alias:        util.ToPtr("bob"),
			},
			args: args{
				id: id.Hex(),
			},
			wantErr: true,
		},
		{
			name: "Request_for_updating_player_is_not_valid_due_to_future_birthdate",
			fields: fields{
				ID:           id.Hex(),
				GovernmentID: "AR-1234567890",
				FirstName:    "Bob",
				MiddleName:   "Sponge",
				LastName:     "Square Pants",
				Category:     nil,
				Birthdate:    util.ToPtr(time.Now().AddDate(20, 0, 0).UTC()),
				PhoneNumber:  "+00 000 000 000",
				Email:        "bobsponge@gmail.com",
				Alias:        util.ToPtr("bob"),
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
		want UpdatePlayer
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
				internalUpdatePlayer: &internalUpdatePlayer{
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
			if got := NewUpdatePlayer(tt.args.dbWriter, tt.args.dbReader); !reflect.DeepEqual(got, tt.want) {
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
		internalUpdatePlayerUsecases *internalUpdatePlayer
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
				internalUpdatePlayerUsecases: &internalUpdatePlayer{
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
					Alias:        util.ToPtr("bob"),
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
					Alias:        util.ToPtr("bob"),
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
					Alias:        util.ToPtr("bob"),
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
					Alias:        util.ToPtr("bob"),
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
				Alias:        util.ToPtr("bob"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareUsecase()
			uc := &updatePlayerUsecase{
				internalUpdatePlayer: tt.fields.internalUpdatePlayerUsecases,
				DBWriter:             tt.fields.DBWriter,
				DBReader:             tt.fields.DBReader,
			}
			got, err := uc.Do(tt.args.ctx, tt.args.id, tt.args.request)
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

func Test_updatePlayerUsecase_UpdatePlayer_Failure(t *testing.T) {
	id := primitive.NewObjectID()
	dbReader := database.NewMockDBReader(gomock.NewController(t))
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))

	type fields struct {
		internalUpdatePlayerUsecases *internalUpdatePlayer
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
			name: "Cannot_update_player_due_to_validation_failure_(missing_government_id)",
			fields: fields{
				internalUpdatePlayerUsecases: &internalUpdatePlayer{
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
					GovernmentID: "",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					Category:     nil,
					Birthdate:    nil,
					PhoneNumber:  "+54 000 000 000",
					Email:        "bobsponge@test.com",
					Alias:        util.ToPtr("bob"),
				},
			},
			prepareUsecase: func() {},
			want:           nil,
			wantErr:        true,
		},
		{
			name: "Cannot_update_player_due_to_government_id_is_not_available",
			fields: fields{
				internalUpdatePlayerUsecases: &internalUpdatePlayer{
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
					Alias:        util.ToPtr("bob"),
				},
			},
			prepareUsecase: func() {
				dbReader.EXPECT().IsAvailable(gomock.Any(), "government_id", gomock.Any()).Return(false, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Cannot_update_player_due_to_government_id_might_not_available",
			fields: fields{
				internalUpdatePlayerUsecases: &internalUpdatePlayer{
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
					Alias:        util.ToPtr("bob"),
				},
			},
			prepareUsecase: func() {
				dbReader.EXPECT().IsAvailable(gomock.Any(), "government_id", gomock.Any()).Return(false, errors.New("could not check government_id availability"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Cannot_update_player_due_to_email_is_not_available",
			fields: fields{
				internalUpdatePlayerUsecases: &internalUpdatePlayer{
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
					Alias:        util.ToPtr("bob"),
				},
			},
			prepareUsecase: func() {
				dbReader.EXPECT().IsAvailable(gomock.Any(), "government_id", "AR-1234567890").Return(true, nil)
				dbReader.EXPECT().IsAvailable(gomock.Any(), "email", "bobsponge@test.com").Return(false, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Cannot_update_player_due_to_email_might_not_available",
			fields: fields{
				internalUpdatePlayerUsecases: &internalUpdatePlayer{
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
					Alias:        util.ToPtr("bob"),
				},
			},
			prepareUsecase: func() {
				dbReader.EXPECT().IsAvailable(gomock.Any(), "government_id", "AR-1234567890").Return(true, nil)
				dbReader.EXPECT().IsAvailable(gomock.Any(), "email", "bobsponge@test.com").Return(false, errors.New("could not check email availability"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Cannot_update_player_due_to_alias_is_not_available",
			fields: fields{
				internalUpdatePlayerUsecases: &internalUpdatePlayer{
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
					Alias:        util.ToPtr("bob"),
				},
			},
			prepareUsecase: func() {
				dbReader.EXPECT().IsAvailable(gomock.Any(), "government_id", "AR-1234567890").Return(true, nil)
				dbReader.EXPECT().IsAvailable(gomock.Any(), "email", "bobsponge@test.com").Return(true, nil)
				dbReader.EXPECT().IsAvailable(gomock.Any(), "alias", "bob").Return(false, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Cannot_update_player_due_to_alias_might_not_available",
			fields: fields{
				internalUpdatePlayerUsecases: &internalUpdatePlayer{
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
					Alias:        util.ToPtr("bob"),
				},
			},
			prepareUsecase: func() {
				dbReader.EXPECT().IsAvailable(gomock.Any(), "government_id", "AR-1234567890").Return(true, nil)
				dbReader.EXPECT().IsAvailable(gomock.Any(), "email", "bobsponge@test.com").Return(true, nil)
				dbReader.EXPECT().IsAvailable(gomock.Any(), "alias", "bob").Return(false, errors.New("could not check alias availability"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Cannot_update_player_due_to_could_not_retrieve_player",
			fields: fields{
				internalUpdatePlayerUsecases: &internalUpdatePlayer{
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
					Alias:        util.ToPtr("bob"),
				},
			},
			prepareUsecase: func() {
				dbReader.EXPECT().IsAvailable(gomock.Any(), "government_id", "AR-1234567890").Return(true, nil)
				dbReader.EXPECT().IsAvailable(gomock.Any(), "email", "bobsponge@test.com").Return(true, nil)
				dbReader.EXPECT().IsAvailable(gomock.Any(), "alias", "bob").Return(true, nil)
				dbReader.EXPECT().GetPlayer(gomock.Any(), gomock.Any()).Return(nil, errors.New("error when retrieving player"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Cannot_update_player_due_to_could_not_update_player",
			fields: fields{
				internalUpdatePlayerUsecases: &internalUpdatePlayer{
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
					Alias:        util.ToPtr("bob"),
				},
			},
			prepareUsecase: func() {
				dbReader.EXPECT().IsAvailable(gomock.Any(), "government_id", "AR-1234567890").Return(true, nil)
				dbReader.EXPECT().IsAvailable(gomock.Any(), "email", "bobsponge@test.com").Return(true, nil)
				dbReader.EXPECT().IsAvailable(gomock.Any(), "alias", "bob").Return(true, nil)
				dbReader.EXPECT().GetPlayer(gomock.Any(), gomock.Any()).Return(&entity.Player{
					ID:           id,
					GovernmentID: "AR-1234567890",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					Category:     nil,
					Birthdate:    nil,
					PhoneNumber:  "+00 000 000 000",
					Email:        "bobsponge@test.com",
					Alias:        util.ToPtr("bob"),
				}, nil)
				dbWriter.EXPECT().UpdatePlayer(gomock.Any(), gomock.Any()).Return(nil, errors.New("could not update player"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareUsecase()
			uc := &updatePlayerUsecase{
				internalUpdatePlayer: tt.fields.internalUpdatePlayerUsecases,
				DBWriter:             tt.fields.DBWriter,
				DBReader:             tt.fields.DBReader,
			}
			got, err := uc.Do(tt.args.ctx, tt.args.id, tt.args.request)
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
