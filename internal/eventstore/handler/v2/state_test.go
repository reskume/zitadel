package handler

import (
	"context"
	"database/sql"
	"database/sql/driver"
	_ "embed"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgconn"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database/mock"
	errs "github.com/zitadel/zitadel/internal/errors"
)

func TestHandler_lockState(t *testing.T) {
	type fields struct {
		projection Projection
		mock       *mock.SQLMock
	}
	type args struct {
		instanceID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		isErr  func(t *testing.T, err error)
	}{
		{
			name: "tx closed",
			fields: fields{
				projection: &projection{
					name: "projection",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExcpectExec(
						lockStateStmt,
						mock.WithExecArgs(
							"projection",
							"instance",
						),
						mock.WithExecErr(sql.ErrTxDone),
					),
				),
			},
			args: args{
				instanceID: "instance",
			},
			isErr: func(t *testing.T, err error) {
				if !errors.Is(err, sql.ErrTxDone) {
					t.Errorf("unexpected error, want: %v got: %v", sql.ErrTxDone, err)
				}
			},
		},
		{
			name: "no rows affeced",
			fields: fields{
				projection: &projection{
					name: "projection",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExcpectExec(
						lockStateStmt,
						mock.WithExecArgs(
							"projection",
							"instance",
						),
						mock.WithExecNoRowsAffected(),
					),
				),
			},
			args: args{
				instanceID: "instance",
			},
			isErr: func(t *testing.T, err error) {
				if !errors.Is(err, errs.ThrowInternal(nil, "V2-lpiK0", "")) {
					t.Errorf("unexpected error: want internal (V2lpiK0), got: %v", err)
				}
			},
		},
		{
			name: "rows affected",
			fields: fields{
				projection: &projection{
					name: "projection",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExcpectExec(
						lockStateStmt,
						mock.WithExecArgs(
							"projection",
							"instance",
						),
						mock.WithExecRowsAffected(1),
					),
				),
			},
			args: args{
				instanceID: "instance",
			},
		},
	}
	for _, tt := range tests {
		if tt.isErr == nil {
			tt.isErr = func(t *testing.T, err error) {
				if err != nil {
					t.Error("expected no error got:", err)
				}
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				projection: tt.fields.projection,
			}

			tx, err := tt.fields.mock.DB.Begin()
			if err != nil {
				t.Fatalf("unable to begin transaction: %v", err)
			}

			err = h.lockState(context.Background(), tx, tt.args.instanceID)
			tt.isErr(t, err)

			tt.fields.mock.Assert(t)
		})
	}
}

func TestHandler_updateLastUpdated(t *testing.T) {
	type fields struct {
		projection Projection
		mock       *mock.SQLMock
	}
	type args struct {
		updatedState *state
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		isErr  func(t *testing.T, err error)
	}{
		{
			name: "update fails",
			fields: fields{
				projection: &projection{
					name: "instance",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExcpectExec(updateStateStmt,
						mock.WithExecErr(sql.ErrTxDone),
					),
				),
			},
			args: args{
				updatedState: &state{
					instanceID:     "instance",
					aggregateType:  "aggregate type",
					aggregateID:    "aggregate id",
					eventTimestamp: time.Now(),
					eventSequence:  42,
				},
			},
			isErr: func(t *testing.T, err error) {
				if !errors.Is(err, sql.ErrTxDone) {
					t.Errorf("unexpected error, want: %v, got %v", sql.ErrTxDone, err)
				}
			},
		},
		{
			name: "no rows affected",
			fields: fields{
				projection: &projection{
					name: "instance",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExcpectExec(updateStateStmt,
						mock.WithExecNoRowsAffected(),
					),
				),
			},
			args: args{
				updatedState: &state{
					instanceID:     "instance",
					aggregateType:  "aggregate type",
					aggregateID:    "aggregate id",
					eventTimestamp: time.Now(),
					eventSequence:  42,
				},
			},
			isErr: func(t *testing.T, err error) {
				if !errors.Is(err, errs.ThrowInternal(nil, "V2-FGEKi", "")) {
					t.Errorf("unexpected error, want: %v, got %v", sql.ErrTxDone, err)
				}
			},
		},
		{
			name: "success",
			fields: fields{
				projection: &projection{
					name: "projection",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExcpectExec(updateStateStmt,
						mock.WithExecArgs(
							"projection",
							"instance",
							mock.AnyType[time.Time]{},
							"aggregate type",
							"aggregate id",
							uint64(42),
						),
						mock.WithExecRowsAffected(1),
					),
				),
			},
			args: args{
				updatedState: &state{
					instanceID:     "instance",
					aggregateType:  "aggregate type",
					aggregateID:    "aggregate id",
					eventTimestamp: time.Now(),
					eventSequence:  42,
				},
			},
		},
	}
	for _, tt := range tests {
		if tt.isErr == nil {
			tt.isErr = func(t *testing.T, err error) {
				if err != nil {
					t.Error("expected no error got:", err)
				}
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			tx, err := tt.fields.mock.DB.Begin()
			if err != nil {
				t.Fatalf("unable to begin transaction: %v", err)
			}

			h := &Handler{
				projection: tt.fields.projection,
			}
			err = h.setState(context.Background(), tx, tt.args.updatedState)

			tt.isErr(t, err)
			tt.fields.mock.Assert(t)
		})
	}
}

func TestHandler_currentState(t *testing.T) {
	testTime := time.Now()
	type fields struct {
		projection Projection
		mock       *mock.SQLMock
	}
	type args struct {
		ctx context.Context
	}
	type want struct {
		currentState *state
		shouldSkip   bool
		isErr        func(t *testing.T, err error)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "connection done",
			fields: fields{
				projection: &projection{
					name: "projection",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExpectQuery(currentStateStmt,
						mock.WithQueryArgs(
							"instance",
							"projection",
						),
						mock.WithQueryErr(sql.ErrConnDone),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
			},
			want: want{
				isErr: func(t *testing.T, err error) {
					if !errors.Is(err, sql.ErrConnDone) {
						t.Errorf("unexpected error, want: %v, got: %v", sql.ErrConnDone, err)
					}
				},
			},
		},
		{
			name: "no row but lock err",
			fields: fields{
				projection: &projection{
					name: "projection",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExpectQuery(currentStateStmt,
						mock.WithQueryArgs(
							"instance",
							"projection",
						),
						mock.WithQueryErr(sql.ErrNoRows),
					),
					mock.ExcpectExec(lockStateStmt,
						mock.WithExecArgs(
							"projection",
							"instance",
						),
						mock.WithExecErr(sql.ErrTxDone),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
			},
			want: want{
				isErr: func(t *testing.T, err error) {
					if !errors.Is(err, sql.ErrTxDone) {
						t.Errorf("unexpected error, want: %v, got: %v", sql.ErrTxDone, err)
					}
				},
			},
		},
		{
			name: "state locked",
			fields: fields{
				projection: &projection{
					name: "projection",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExpectQuery(currentStateStmt,
						mock.WithQueryArgs(
							"instance",
							"projection",
						),
						mock.WithQueryErr(&pgconn.PgError{Code: "55P03"}),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
			},
			want: want{
				shouldSkip: true,
			},
		},
		{
			name: "success",
			fields: fields{
				projection: &projection{
					name: "projection",
				},
				mock: mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExpectQuery(currentStateStmt,
						mock.WithQueryArgs(
							"instance",
							"projection",
						),
						mock.WithQueryResult(
							[]string{"event_date", "aggregate_type", "aggregate_id", "event_sequence"},
							[][]driver.Value{
								{
									testTime,
									"aggregate type",
									"aggregate id",
									int64(42),
								},
							},
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance"),
			},
			want: want{
				currentState: &state{
					instanceID:     "instance",
					aggregateType:  "aggregate type",
					aggregateID:    "aggregate id",
					eventTimestamp: testTime,
					eventSequence:  42,
				},
			},
		},
	}
	for _, tt := range tests {
		if tt.want.isErr == nil {
			tt.want.isErr = func(t *testing.T, err error) {
				if err != nil {
					t.Error("expected no error got:", err)
				}
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				projection: tt.fields.projection,
			}

			tx, err := tt.fields.mock.DB.Begin()
			if err != nil {
				t.Fatalf("unable to begin transaction: %v", err)
			}

			gotCurrentState, gotShouldSkip, err := h.currentState(tt.args.ctx, tx)

			tt.want.isErr(t, err)
			if !reflect.DeepEqual(gotCurrentState, tt.want.currentState) {
				t.Errorf("Handler.currentState() gotCurrentState = %v, want %v", gotCurrentState, tt.want.currentState)
			}
			if gotShouldSkip != tt.want.shouldSkip {
				t.Errorf("Handler.currentState() gotShouldSkip = %v, want %v", gotShouldSkip, tt.want.shouldSkip)
			}
			tt.fields.mock.Assert(t)
		})
	}
}