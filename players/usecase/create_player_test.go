package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/Neniel/gotennis/lib/database"
	"github.com/Neniel/gotennis/lib/entity"
	"github.com/Neniel/gotennis/lib/util"
	"go.uber.org/mock/gomock"
)

func Test_createPlayerUsecase_CreatePlayer_Success(t *testing.T) {
	dbReaderMock := database.NewMockDBReader(gomock.NewController(t))
	dbWriterMock := database.NewMockDBWriter(gomock.NewController(t))

	type fields struct {
		internalCreatePlayerUsecases *internalCreatePlayer
		DBWriter                     database.DBWriter
	}
	type args struct {
		ctx     context.Context
		request *CreatePlayerRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		prepareMocks func()
		want         *entity.Player
		wantErr      bool
	}{
		{
			name: "Try_to_create_player_with_all_valid_data",
			fields: fields{
				internalCreatePlayerUsecases: &internalCreatePlayer{
					ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReaderMock),
					ValidateEmail:        NewValidateEmailUsecase(dbReaderMock),
					ValidateAlias:        NewValidateAliasUsecase(dbReaderMock),
				},
				DBWriter: dbWriterMock,
			},
			args: args{
				ctx: context.Background(),
				request: &CreatePlayerRequest{
					GovernmentID: "1234567890",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					Category:     nil,
					Birthdate:    nil,
					PhoneNumber:  "+54 000 000 000",
					Email:        "spongebob@test.com",
					Alias:        util.ToPtr(""),
				},
			},
			prepareMocks: func() {
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "government_id", "1234567890").Return(true, nil)
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "email", "spongebob@test.com").Return(true, nil)
				dbWriterMock.EXPECT().AddPlayer(gomock.Any(), gomock.Any()).Return(&entity.Player{
					GovernmentID: "1234567890",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					PhoneNumber:  "+54 000 000 000",
					Email:        "spongebob@test.com",
					Alias:        util.ToPtr(""),
				}, nil)
			},
			want: &entity.Player{
				GovernmentID: "1234567890",
				FirstName:    "Bob",
				MiddleName:   "Sponge",
				LastName:     "Square Pants",
				PhoneNumber:  "+54 000 000 000",
				Email:        "spongebob@test.com",
				Alias:        util.ToPtr(""),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks()
			uc := &createPlayer{
				internalCreatePlayer: tt.fields.internalCreatePlayerUsecases,
				DBWriter:             tt.fields.DBWriter,
			}
			got, err := uc.Do(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("createPlayerUsecase.CreatePlayer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createPlayerUsecase.CreatePlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createPlayerUsecase_CreatePlayer_Failure(t *testing.T) {
	dbReaderMock := database.NewMockDBReader(gomock.NewController(t))
	dbWriterMock := database.NewMockDBWriter(gomock.NewController(t))

	type fields struct {
		internalCreatePlayerUsecases *internalCreatePlayer
		DBWriter                     database.DBWriter
	}
	type args struct {
		ctx     context.Context
		request *CreatePlayerRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		prepareMocks func()
		want         *entity.Player
		wantErr      bool
	}{
		{
			name: "Fails_when_GovernmentID_in_request_is_empty",
			fields: fields{
				internalCreatePlayerUsecases: &internalCreatePlayer{
					ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReaderMock),
					ValidateEmail:        NewValidateEmailUsecase(dbReaderMock),
					ValidateAlias:        NewValidateAliasUsecase(dbReaderMock),
				},
				DBWriter: dbWriterMock,
			},
			args: args{
				ctx: context.Background(),
				request: &CreatePlayerRequest{
					GovernmentID: "",
				},
			},
			prepareMocks: func() {},
			want:         nil,
			wantErr:      true,
		},
		{
			name: "Fails_when_Email_in_request_is_empty",
			fields: fields{
				internalCreatePlayerUsecases: &internalCreatePlayer{
					ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReaderMock),
					ValidateEmail:        NewValidateEmailUsecase(dbReaderMock),
					ValidateAlias:        NewValidateAliasUsecase(dbReaderMock),
				},
				DBWriter: dbWriterMock,
			},
			args: args{
				ctx: context.Background(),
				request: &CreatePlayerRequest{
					GovernmentID: "1234567890",
					Email:        "",
				},
			},
			prepareMocks: func() {},
			want:         nil,
			wantErr:      true,
		},
		{
			name: "Fails_when_FirstName_in_request_is_empty",
			fields: fields{
				internalCreatePlayerUsecases: &internalCreatePlayer{
					ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReaderMock),
					ValidateEmail:        NewValidateEmailUsecase(dbReaderMock),
					ValidateAlias:        NewValidateAliasUsecase(dbReaderMock),
				},
				DBWriter: dbWriterMock,
			},
			args: args{
				ctx: context.Background(),
				request: &CreatePlayerRequest{
					GovernmentID: "1234567890",
					Email:        "1234567890@test.com",
					FirstName:    "",
				},
			},
			prepareMocks: func() {},
			want:         nil,
			wantErr:      true,
		},
		{
			name: "Fails_when_LastName_in_request_is_empty",
			fields: fields{
				internalCreatePlayerUsecases: &internalCreatePlayer{
					ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReaderMock),
					ValidateEmail:        NewValidateEmailUsecase(dbReaderMock),
					ValidateAlias:        NewValidateAliasUsecase(dbReaderMock),
				},
				DBWriter: dbWriterMock,
			},
			args: args{
				ctx: context.Background(),
				request: &CreatePlayerRequest{
					GovernmentID: "1234567890",
					Email:        "1234567890@test.com",
					FirstName:    "Bob",
					LastName:     "",
				},
			},
			prepareMocks: func() {},
			want:         nil,
			wantErr:      true,
		},
		{
			name: "Fails_when_the_provided_GovernmentID_is_not_available",
			fields: fields{
				internalCreatePlayerUsecases: &internalCreatePlayer{
					ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReaderMock),
					ValidateEmail:        NewValidateEmailUsecase(dbReaderMock),
					ValidateAlias:        NewValidateAliasUsecase(dbReaderMock),
				},
				DBWriter: dbWriterMock,
			},
			args: args{
				ctx: context.Background(),
				request: &CreatePlayerRequest{
					GovernmentID: "1234567890",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					Category:     nil,
					Birthdate:    nil,
					PhoneNumber:  "+54 000 000 000",
					Email:        "spongebob@test.com",
					Alias:        util.ToPtr(""),
				},
			},
			prepareMocks: func() {
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "government_id", "1234567890").Return(false, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Fails_when_the_provided_GovernmentID_might_be_available_but_could_not_check_availability",
			fields: fields{
				internalCreatePlayerUsecases: &internalCreatePlayer{
					ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReaderMock),
					ValidateEmail:        NewValidateEmailUsecase(dbReaderMock),
					ValidateAlias:        NewValidateAliasUsecase(dbReaderMock),
				},
				DBWriter: dbWriterMock,
			},
			args: args{
				ctx: context.Background(),
				request: &CreatePlayerRequest{
					GovernmentID: "1234567890",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					Category:     nil,
					Birthdate:    nil,
					PhoneNumber:  "+54 000 000 000",
					Email:        "spongebob@test.com",
					Alias:        util.ToPtr(""),
				},
			},
			prepareMocks: func() {
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "government_id", "1234567890").Return(false, errors.New("error checking government_id availability"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Fails_when_the_provided_Email_is_not_available",
			fields: fields{
				internalCreatePlayerUsecases: &internalCreatePlayer{
					ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReaderMock),
					ValidateEmail:        NewValidateEmailUsecase(dbReaderMock),
					ValidateAlias:        NewValidateAliasUsecase(dbReaderMock),
				},
				DBWriter: dbWriterMock,
			},
			args: args{
				ctx: context.Background(),
				request: &CreatePlayerRequest{
					GovernmentID: "1234567890",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					Category:     nil,
					Birthdate:    nil,
					PhoneNumber:  "+54 000 000 000",
					Email:        "spongebob@test.com",
					Alias:        util.ToPtr(""),
				},
			},
			prepareMocks: func() {
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "government_id", "1234567890").Return(true, nil)
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "email", "spongebob@test.com").Return(false, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Fails_when_the_provided_Email_might_be_available_but_could_not_check_availability",
			fields: fields{
				internalCreatePlayerUsecases: &internalCreatePlayer{
					ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReaderMock),
					ValidateEmail:        NewValidateEmailUsecase(dbReaderMock),
					ValidateAlias:        NewValidateAliasUsecase(dbReaderMock),
				},
				DBWriter: dbWriterMock,
			},
			args: args{
				ctx: context.Background(),
				request: &CreatePlayerRequest{
					GovernmentID: "1234567890",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					Category:     nil,
					Birthdate:    nil,
					PhoneNumber:  "+54 000 000 000",
					Email:        "spongebob@test.com",
					Alias:        util.ToPtr(""),
				},
			},
			prepareMocks: func() {
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "government_id", "1234567890").Return(true, nil)
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "email", "spongebob@test.com").Return(false, errors.New("could not check email availability"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Fails_when_the_provided_Alias_is_not_available",
			fields: fields{
				internalCreatePlayerUsecases: &internalCreatePlayer{
					ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReaderMock),
					ValidateEmail:        NewValidateEmailUsecase(dbReaderMock),
					ValidateAlias:        NewValidateAliasUsecase(dbReaderMock),
				},
				DBWriter: dbWriterMock,
			},
			args: args{
				ctx: context.Background(),
				request: &CreatePlayerRequest{
					GovernmentID: "1234567890",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					Category:     nil,
					Birthdate:    nil,
					PhoneNumber:  "+54 000 000 000",
					Email:        "spongebob@test.com",
					Alias:        util.ToPtr("bob"),
				},
			},
			prepareMocks: func() {
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "government_id", "1234567890").Return(true, nil)
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "email", "spongebob@test.com").Return(true, nil)
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "alias", "bob").Return(false, nil)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Fails_when_the_provided_Alias_might_be_available_but_could_not_check_availability",
			fields: fields{
				internalCreatePlayerUsecases: &internalCreatePlayer{
					ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReaderMock),
					ValidateEmail:        NewValidateEmailUsecase(dbReaderMock),
					ValidateAlias:        NewValidateAliasUsecase(dbReaderMock),
				},
				DBWriter: dbWriterMock,
			},
			args: args{
				ctx: context.Background(),
				request: &CreatePlayerRequest{
					GovernmentID: "1234567890",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					Category:     nil,
					Birthdate:    nil,
					PhoneNumber:  "+54 000 000 000",
					Email:        "spongebob@test.com",
					Alias:        util.ToPtr("bob"),
				},
			},
			prepareMocks: func() {
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "government_id", "1234567890").Return(true, nil)
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "email", "spongebob@test.com").Return(true, nil)
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "alias", "bob").Return(false, errors.New("could not check email availability"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Fails_when_all_data_is_correct_but_cannot_save_player_into_the_database",
			fields: fields{
				internalCreatePlayerUsecases: &internalCreatePlayer{
					ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReaderMock),
					ValidateEmail:        NewValidateEmailUsecase(dbReaderMock),
					ValidateAlias:        NewValidateAliasUsecase(dbReaderMock),
				},
				DBWriter: dbWriterMock,
			},
			args: args{
				ctx: context.Background(),
				request: &CreatePlayerRequest{
					GovernmentID: "1234567890",
					FirstName:    "Bob",
					MiddleName:   "Sponge",
					LastName:     "Square Pants",
					Category:     nil,
					Birthdate:    nil,
					PhoneNumber:  "+54 000 000 000",
					Email:        "spongebob@test.com",
					Alias:        util.ToPtr("bob"),
				},
			},
			prepareMocks: func() {
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "government_id", "1234567890").Return(true, nil)
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "email", "spongebob@test.com").Return(true, nil)
				dbReaderMock.EXPECT().IsAvailable(gomock.Any(), "alias", "bob").Return(true, nil)
				dbWriterMock.EXPECT().AddPlayer(gomock.Any(), gomock.Any()).Return(nil, errors.New("error saving player into the database"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareMocks()
			uc := &createPlayer{
				internalCreatePlayer: tt.fields.internalCreatePlayerUsecases,
				DBWriter:             tt.fields.DBWriter,
			}
			got, err := uc.Do(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("createPlayerUsecase.CreatePlayer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createPlayerUsecase.CreatePlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCreatePlayerUsecase(t *testing.T) {
	dbWriter := database.NewMockDBWriter(gomock.NewController(t))
	dbReader := database.NewMockDBReader(gomock.NewController(t))

	type args struct {
		dbWriter database.DBWriter
		dbReader database.DBReader
	}
	tests := []struct {
		name string
		args args
		want CreatePlayer
	}{
		{
			name: "Should_get_a_create_player_usecase",
			args: args{
				dbWriter: dbWriter,
				dbReader: dbReader,
			},
			want: &createPlayer{
				internalCreatePlayer: &internalCreatePlayer{
					ValidateGovernmentID: NewValidateGovernmentIDUsecase(dbReader),
					ValidateEmail:        NewValidateEmailUsecase(dbReader),
					ValidateAlias:        NewValidateAliasUsecase(dbReader),
				},
				DBWriter: dbWriter,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCreatePlayer(tt.args.dbWriter, tt.args.dbReader); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCreatePlayerUsecase() = %v, want %v", got, tt.want)
			}
		})
	}
}
