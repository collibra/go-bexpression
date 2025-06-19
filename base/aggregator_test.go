package base

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/collibra/go-bexpression/utils"
)

func TestAggregator_Validate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}
	type testCase[T Comparison] struct {
		name    string
		a       Aggregator[T]
		args    args
		wantErr bool
	}
	tests := []testCase[*MockComparison]{
		{
			name: "no operands",
			args: args{ctx: context.Background()},
			a: Aggregator[*MockComparison]{
				Operator: AggregatorOperatorAnd,
				Operands: nil,
			},
			wantErr: true,
		},
		{
			name: "valid operands",
			args: args{ctx: context.Background()},
			a: Aggregator[*MockComparison]{
				Operator: AggregatorOperatorOr,
				Operands: []BinaryExpression[*MockComparison]{
					{
						Literal: utils.Ptr(true),
					},
					{
						Literal: utils.Ptr(false),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid operands",
			args: args{ctx: context.Background()},
			a: Aggregator[*MockComparison]{
				Operator: AggregatorOperatorOr,
				Operands: []BinaryExpression[*MockComparison]{
					{
						Literal: utils.Ptr(true),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.a.Validate(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAggregator_Accept(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		mockSetup func(mockVisitor *MockVisitor)
	}
	type testCase[T Comparison] struct {
		name    string
		a       Aggregator[T]
		args    args
		wantErr bool
	}
	tests := []testCase[*MockComparison]{
		{
			name: "happy path",
			args: args{
				ctx: context.Background(),
				mockSetup: func(visitor *MockVisitor) {
					aggregatorEnterCall := visitor.EXPECT().EnterExpressionElement(mock.Anything, mock.AnythingOfType("*base.Aggregator[*github.com/collibra/go-bexpression/base.MockComparison]")).Return(nil).Once()
					beEnterCall := visitor.EXPECT().EnterExpressionElement(mock.Anything, mock.AnythingOfType("*base.BinaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return(nil).Twice().NotBefore(aggregatorEnterCall)

					literalTrueCall := visitor.EXPECT().Literal(mock.Anything, true).Return(nil).Once()

					visitor.EXPECT().Literal(mock.Anything, AggregatorOperatorAnd).Once().Return(nil).NotBefore(literalTrueCall)

					literalFalseCall := visitor.EXPECT().Literal(mock.Anything, false).Return(nil).Once()

					beLeaveCall := visitor.EXPECT().LeaveExpressionElement(mock.Anything, mock.AnythingOfType("*base.BinaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return().Twice()

					visitor.EXPECT().LeaveExpressionElement(mock.Anything, mock.AnythingOfType("*base.Aggregator[*github.com/collibra/go-bexpression/base.MockComparison]")).Return().NotBefore(beLeaveCall, literalTrueCall, literalFalseCall, beEnterCall, aggregatorEnterCall)

				},
			},
			a: Aggregator[*MockComparison]{
				Operator: AggregatorOperatorAnd,
				Operands: []BinaryExpression[*MockComparison]{
					{
						Literal: utils.Ptr(true),
					},
					{
						Literal: utils.Ptr(false),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visitor := NewMockVisitor(t)
			tt.args.mockSetup(visitor)

			if err := tt.a.Accept(tt.args.ctx, visitor); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
