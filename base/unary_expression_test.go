package base

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/collibra/go-bexpression/utils"
)

func TestUnaryExpression_Accept(t *testing.T) {
	type args struct {
		ctx   context.Context
		setup func(visitor *MockVisitor)
	}
	type testCase[T Comparison] struct {
		name    string
		n       UnaryExpression[T]
		args    args
		wantErr bool
	}
	tests := []testCase[*MockComparison]{
		{
			name: "happy path",
			n: UnaryExpression[*MockComparison]{
				Operator: UnaryOperatorNot,
				Operand: BinaryExpression[*MockComparison]{
					Literal: utils.Ptr(true),
				},
			},
			args: args{
				ctx: context.Background(),
				setup: func(visitor *MockVisitor) {
					enterUnaryCall := visitor.EXPECT().EnterExpressionElement(mock.Anything, mock.AnythingOfType("*base.UnaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return(nil).Once()
					literalOperatorCall := visitor.EXPECT().Literal(mock.Anything, UnaryOperatorNot).Return(nil).Once().NotBefore(enterUnaryCall)

					enterBexpressionCall := visitor.EXPECT().EnterExpressionElement(mock.Anything, mock.AnythingOfType("*base.BinaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return(nil).Once().NotBefore(literalOperatorCall)
					literalCall := visitor.EXPECT().Literal(mock.Anything, true).Return(nil).Once().NotBefore(enterBexpressionCall)

					leaveBexpressionCall := visitor.EXPECT().LeaveExpressionElement(mock.Anything, mock.AnythingOfType("*base.BinaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return().Once().NotBefore(literalCall)
					visitor.EXPECT().LeaveExpressionElement(mock.Anything, mock.AnythingOfType("*base.UnaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return().NotBefore(leaveBexpressionCall)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVisitor := NewMockVisitor(t)
			tt.args.setup(mockVisitor)

			if err := tt.n.Accept(tt.args.ctx, mockVisitor); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnaryExpression_Validate(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type testCase[T Comparison] struct {
		name    string
		n       UnaryExpression[T]
		args    args
		wantErr bool
	}
	tests := []testCase[*MockComparison]{
		{
			name: "valid bexpr",
			n: UnaryExpression[*MockComparison]{
				Operator: UnaryOperatorNot,
				Operand: BinaryExpression[*MockComparison]{
					Literal: utils.Ptr(true),
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "invalid bexpr",
			n: UnaryExpression[*MockComparison]{
				Operator: UnaryOperatorNot,
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Validate(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
