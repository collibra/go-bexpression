package datacomparison

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/collibra/go-bexpression/base"
	"github.com/collibra/go-bexpression/utils"
)

func TestOperand_Accept(t *testing.T) {
	type fields struct {
		Reference *Reference
		Literal   *Literal
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
			name: "reference",
			fields: fields{
				Reference: &Reference{
					EntityID: "entity_id",
				},
			},
			args: args{
				ctx: context.Background(),
				visitorSetup: func(visitor *base.MockVisitor) {
					enterCall := visitor.EXPECT().EnterExpressionElement(mock.Anything, mock.AnythingOfType("*datacomparison.Operand")).Return(nil).Once()
					refListeralCall := visitor.EXPECT().Literal(mock.Anything, &Reference{EntityID: "entity_id"}).Return(nil).Once().NotBefore(enterCall)
					visitor.EXPECT().LeaveExpressionElement(mock.Anything, mock.AnythingOfType("*datacomparison.Operand")).Return().Once().NotBefore(refListeralCall)

				},
			},
			wantErr: false,
		},
		{
			name: "literal",
			fields: fields{
				Literal: &Literal{
					Str: utils.Ptr("str"),
				},
			},
			args: args{
				ctx: context.Background(),
				visitorSetup: func(visitor *base.MockVisitor) {
					enterCall := visitor.EXPECT().EnterExpressionElement(mock.Anything, mock.AnythingOfType("*datacomparison.Operand")).Return(nil).Once()
					literalCall := visitor.EXPECT().Literal(mock.Anything, "str").Return(nil).Once().NotBefore(enterCall)
					visitor.EXPECT().LeaveExpressionElement(mock.Anything, mock.AnythingOfType("*datacomparison.Operand")).Return().Once().NotBefore(literalCall)

				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visitor := base.NewMockVisitor(t)
			tt.args.visitorSetup(visitor)

			o := &Operand{
				Reference: tt.fields.Reference,
				Literal:   tt.fields.Literal,
			}
			if err := o.Accept(tt.args.ctx, visitor); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOperand_ToGql(t *testing.T) {
	type fields struct {
		Reference *Reference
		Literal   *Literal
	}
	tests := []struct {
		name    string
		fields  fields
		want    DataComparisonOperand
		wantErr bool
	}{
		{
			name: "reference",
			fields: fields{
				Reference: &Reference{
					EntityID: "entity_id",
				},
			},
			want: &Reference{
				EntityID: "entity_id",
			},
			wantErr: false,
		},
		{
			name: "literal",
			fields: fields{
				Literal: &Literal{
					Int: utils.Ptr(1),
				},
			},
			want: &LiteralInt{
				Value: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Operand{
				Reference: tt.fields.Reference,
				Literal:   tt.fields.Literal,
			}
			got, err := o.ToGql()
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

func TestOperand_Validate(t *testing.T) {
	type fields struct {
		Reference *Reference
		Literal   *Literal
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "invalid operand",
			fields: fields{},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "valid reference",
			fields: fields{
				Reference: &Reference{},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "valid literal",
			fields: fields{
				Literal: &Literal{
					Int: utils.Ptr(1),
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "invalid literal",
			fields: fields{
				Literal: &Literal{},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Operand{
				Reference: tt.fields.Reference,
				Literal:   tt.fields.Literal,
			}
			if err := o.Validate(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
