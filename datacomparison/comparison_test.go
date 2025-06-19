package datacomparison

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/collibra/go-bexpression/base"
	"github.com/collibra/go-bexpression/utils"
)

func TestDataComparison_Accept(t *testing.T) {
	type fields struct {
		Operator     ComparisonOperator
		LeftOperand  Operand
		RightOperand Operand
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
			name: "happy path",
			fields: fields{
				Operator:     ComparisonOperatorEqual,
				LeftOperand:  Operand{Literal: &Literal{Int: utils.Ptr(1)}},
				RightOperand: Operand{Literal: &Literal{Int: utils.Ptr(5)}},
			},
			args: args{
				ctx: context.Background(),
				visitorSetup: func(visitor *base.MockVisitor) {
					enterCall := visitor.EXPECT().EnterExpressionElement(mock.Anything, mock.AnythingOfType("*datacomparison.DataComparison")).Return(nil).Once()
					visitor.EXPECT().EnterExpressionElement(mock.Anything, mock.AnythingOfType("*datacomparison.Operand")).Return(nil).Twice().NotBefore(enterCall)
					literalCall1 := visitor.EXPECT().Literal(mock.Anything, 1).Return(nil).Once()
					literalOperatorCall := visitor.EXPECT().Literal(mock.Anything, ComparisonOperatorEqual).Return(nil).Once().NotBefore(literalCall1)
					visitor.EXPECT().Literal(mock.Anything, 5).Return(nil).Once().NotBefore(literalOperatorCall)
					visitor.EXPECT().LeaveExpressionElement(mock.Anything, mock.AnythingOfType("*datacomparison.Operand")).Return().Twice()
					visitor.EXPECT().LeaveExpressionElement(mock.Anything, mock.AnythingOfType("*datacomparison.DataComparison")).Return()
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visitor := base.NewMockVisitor(t)
			tt.args.visitorSetup(visitor)

			d := &DataComparison{
				Operator:     tt.fields.Operator,
				LeftOperand:  tt.fields.LeftOperand,
				RightOperand: tt.fields.RightOperand,
			}
			if err := d.Accept(tt.args.ctx, visitor); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDataComparison_Validate(t *testing.T) {
	type fields struct {
		Operator     ComparisonOperator
		LeftOperand  Operand
		RightOperand Operand
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
			name: "valid operands",
			fields: fields{
				Operator:     ComparisonOperatorEqual,
				LeftOperand:  Operand{Literal: &Literal{Int: utils.Ptr(1)}},
				RightOperand: Operand{Literal: &Literal{Int: utils.Ptr(5)}},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "invalid left operand",
			fields: fields{
				Operator:     ComparisonOperatorEqual,
				LeftOperand:  Operand{},
				RightOperand: Operand{Literal: &Literal{Int: utils.Ptr(5)}},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "invalid right operand",
			fields: fields{
				Operator:     ComparisonOperatorEqual,
				LeftOperand:  Operand{Literal: &Literal{Int: utils.Ptr(1)}},
				RightOperand: Operand{},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DataComparison{
				Operator:     tt.fields.Operator,
				LeftOperand:  tt.fields.LeftOperand,
				RightOperand: tt.fields.RightOperand,
			}
			if err := d.Validate(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
