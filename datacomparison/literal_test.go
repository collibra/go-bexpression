package datacomparison

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/collibra/go-bexpression/base"
	"github.com/collibra/go-bexpression/utils"
)

func TestLiteral_Accept(t *testing.T) {
	tNow := time.Now()

	type fields struct {
		Bool      *bool
		Int       *int
		Float     *float64
		Str       *string
		Timestamp *time.Time
	}
	type args struct {
		ctx          context.Context
		visitorSetup func(visitor *base.MockVisitor)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "bool",
			fields: fields{
				Bool: utils.Ptr(true),
			},
			args: args{
				ctx: context.Background(),
				visitorSetup: func(visitor *base.MockVisitor) {
					visitor.EXPECT().Literal(mock.Anything, true).Return(nil).Once()
				},
			},
			wantErr: false,
		},
		{
			name: "int",
			fields: fields{
				Int: utils.Ptr(1),
			},
			args: args{
				ctx: context.Background(),
				visitorSetup: func(visitor *base.MockVisitor) {
					visitor.EXPECT().Literal(mock.Anything, 1).Return(nil).Once()
				},
			},
		},
		{
			name: "float",
			fields: fields{
				Float: utils.Ptr(1.1),
			},
			args: args{
				ctx: context.Background(),
				visitorSetup: func(visitor *base.MockVisitor) {
					visitor.EXPECT().Literal(mock.Anything, 1.1).Return(nil).Once()
				},
			},
		},
		{
			name: "string",
			fields: fields{
				Str: utils.Ptr("str"),
			},
			args: args{
				ctx: context.Background(),
				visitorSetup: func(visitor *base.MockVisitor) {
					visitor.EXPECT().Literal(mock.Anything, "str").Return(nil).Once()
				},
			},
		},
		{
			name: "timestamp",
			fields: fields{
				Timestamp: utils.Ptr(tNow),
			},
			args: args{
				ctx: context.Background(),
				visitorSetup: func(visitor *base.MockVisitor) {
					visitor.EXPECT().Literal(mock.Anything, tNow).Return(nil).Once()
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visitorMock := base.NewMockVisitor(t)
			tt.args.visitorSetup(visitorMock)

			l := &Literal{
				Bool:      tt.fields.Bool,
				Int:       tt.fields.Int,
				Float:     tt.fields.Float,
				Str:       tt.fields.Str,
				Timestamp: tt.fields.Timestamp,
			}
			if err := l.Accept(tt.args.ctx, visitorMock); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLiteral_ToGql(t *testing.T) {
	tNow := time.Now()

	type fields struct {
		Bool      *bool
		Int       *int
		Float     *float64
		Str       *string
		Timestamp *time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    DataComparisonOperand
		wantErr bool
	}{
		{
			name: "bool",
			fields: fields{
				Bool:      utils.Ptr(true),
				Int:       nil,
				Float:     nil,
				Str:       nil,
				Timestamp: nil,
			},
			want:    &LiteralBool{Value: true},
			wantErr: false,
		},
		{
			name: "int",
			fields: fields{
				Bool:      nil,
				Int:       utils.Ptr(1),
				Float:     nil,
				Str:       nil,
				Timestamp: nil,
			},
			want:    &LiteralInt{Value: 1},
			wantErr: false,
		},
		{
			name: "float",
			fields: fields{
				Bool:      nil,
				Int:       nil,
				Float:     utils.Ptr(1.1),
				Str:       nil,
				Timestamp: nil,
			},
			want:    &LiteralFloat{Value: 1.1},
			wantErr: false,
		},
		{
			name: "string",
			fields: fields{
				Bool:      nil,
				Int:       nil,
				Float:     nil,
				Str:       utils.Ptr("test"),
				Timestamp: nil,
			},
			want:    &LiteralString{Value: "test"},
			wantErr: false,
		},
		{
			name: "timestamp",
			fields: fields{
				Bool:      nil,
				Int:       nil,
				Float:     nil,
				Str:       nil,
				Timestamp: utils.Ptr(tNow),
			},
			want:    &LiteralTime{Value: tNow},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Literal{
				Bool:      tt.fields.Bool,
				Int:       tt.fields.Int,
				Float:     tt.fields.Float,
				Str:       tt.fields.Str,
				Timestamp: tt.fields.Timestamp,
			}
			got, err := l.ToGql()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToGql() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToGql() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLiteral_Validate(t *testing.T) {
	type fields struct {
		Bool      *bool
		Int       *int
		Float     *float64
		Str       *string
		Timestamp *time.Time
	}
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "no literals defined",
			fields: fields{},
			args: args{
				in0: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "one literal defined",
			fields: fields{
				Bool: utils.Ptr(true),
			},
			wantErr: false,
		},
		{
			name: "multiple literals defined",
			fields: fields{
				Bool:  utils.Ptr(true),
				Float: utils.Ptr(1.1),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Literal{
				Bool:      tt.fields.Bool,
				Int:       tt.fields.Int,
				Float:     tt.fields.Float,
				Str:       tt.fields.Str,
				Timestamp: tt.fields.Timestamp,
			}
			if err := l.Validate(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
