package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/inasknh/simple-poke-app/internal/model"
	"reflect"
	"regexp"
	"testing"
)

func Test_repository_CreateBerry(t *testing.T) {
	type args struct {
		ctx     context.Context
		berries []model.Berry
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		mockCall func(mock sqlmock.Sqlmock)
	}{
		{
			name: "given empty berries in request should return nil error",
			args: args{
				ctx:     context.Background(),
				berries: nil,
			},
			wantErr: false,
			mockCall: func(mock sqlmock.Sqlmock) {

			},
		},
		{
			name: "given an error when execute query should return an error",
			args: args{
				ctx: context.Background(),
				berries: []model.Berry{
					{
						Name: "1",
						URL:  "1-url",
					},
					{
						Name: "2",
						URL:  "2-url",
					},
				},
			},
			wantErr: true,
			mockCall: func(mock sqlmock.Sqlmock) {
				query := "INSERT INTO berries (name, url) VALUES (?, ?),(?, ?)"
				mock.
					ExpectExec(regexp.QuoteMeta(query)).WithArgs(
					"1",
					"1-url",
					"2",
					"2-url",
				).
					WillReturnError(errors.New("error"))
			},
		},
		{
			name: "given no error when execute query should return nil error",
			args: args{
				ctx: context.Background(),
				berries: []model.Berry{
					{
						Name: "1",
						URL:  "1-url",
					},
					{
						Name: "2",
						URL:  "2-url",
					},
				},
			},
			wantErr: true,
			mockCall: func(mock sqlmock.Sqlmock) {
				query := "INSERT INTO berries (name, url) VALUES (?, ?),(?, ?)"
				mock.
					ExpectExec(regexp.QuoteMeta(query)).WithArgs(
					"1",
					"1-url",
					"2",
					"2-url",
				).
					WillReturnResult(sqlmock.NewResult(2, 0))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				_ = db.Close()
			}(db)

			tt.mockCall(mock)

			r := &repository{
				db: db,
			}

			if err := r.CreateBerry(tt.args.ctx, tt.args.berries); (err != nil) != tt.wantErr {
				t.Errorf("CreateBerry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_repository_FetchBerries(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		want     *model.BerriesResponse
		wantErr  bool
		mockCall func(mock sqlmock.Sqlmock)
	}{
		{
			name: "given an error when execute query should return nil and an error",
			args: args{
				ctx: context.Background(),
			},
			want:    nil,
			wantErr: true,
			mockCall: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(getAllBerries).WillReturnError(errors.New("any error"))
			},
		},
		{
			name: "given happy flow should return response and no error",
			args: args{
				ctx: context.Background(),
			},
			want: &model.BerriesResponse{Berries: []model.Berry{
				{
					Name: "1",
					URL:  "1",
				},
			}},
			wantErr: false,
			mockCall: func(mock sqlmock.Sqlmock) {
				columns := []string{
					"name",
					"url",
				}
				mockRes := mock.NewRows(columns).AddRow(
					"1",
					"1",
				)
				mock.ExpectQuery(getAllBerries).WillReturnRows(mockRes)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				_ = db.Close()
			}(db)

			tt.mockCall(mock)

			r := &repository{
				db: db,
			}
			got, err := r.FetchBerries(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchBerries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchBerries() got = %v, want %v", got, tt.want)
			}
		})
	}
}
