package base

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/collibra/go-bexpression/utils"
)

func TestBinaryExpression_Validate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}
	type testCase[T Comparison] struct {
		name    string
		b       func() BinaryExpression[T]
		args    args
		wantErr bool
	}
	tests := []testCase[*MockComparison]{
		{
			name: "valid literal expression",
			b: func() BinaryExpression[*MockComparison] {
				return BinaryExpression[*MockComparison]{
					Literal: utils.Ptr(true),
				}
			},
			wantErr: false,
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "valid comparison expression",
			b: func() BinaryExpression[*MockComparison] {
				compMock := NewMockComparison(t)
				compMock.EXPECT().Validate(mock.Anything).Return(nil).Once()

				return BinaryExpression[*MockComparison]{
					Comparison: compMock,
				}
			},
			wantErr: false,
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "valid aggregator expression",
			b: func() BinaryExpression[*MockComparison] {
				return BinaryExpression[*MockComparison]{
					Aggregator: &Aggregator[*MockComparison]{
						Operator: AggregatorOperatorAnd,
						Operands: []BinaryExpression[*MockComparison]{
							{
								Literal: utils.Ptr(true),
							},
						},
					},
				}
			},
			wantErr: false,
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "valid unary expression",
			b: func() BinaryExpression[*MockComparison] {
				return BinaryExpression[*MockComparison]{
					UnaryExpression: &UnaryExpression[*MockComparison]{
						Operator: UnaryOperatorNot,
						Operand: BinaryExpression[*MockComparison]{
							Literal: utils.Ptr(true),
						},
					},
				}
			},
			wantErr: false,
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "invalid binary expression",
			b: func() BinaryExpression[*MockComparison] {
				return BinaryExpression[*MockComparison]{}
			},
			wantErr: true,
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "invalid comparison expression",
			b: func() BinaryExpression[*MockComparison] {
				compMock := NewMockComparison(t)
				compMock.EXPECT().Validate(mock.Anything).Return(errors.New("boom")).Once()

				return BinaryExpression[*MockComparison]{
					Comparison: compMock,
				}
			},
			wantErr: true,
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "invalid aggregator expression",
			b: func() BinaryExpression[*MockComparison] {
				return BinaryExpression[*MockComparison]{
					Aggregator: &Aggregator[*MockComparison]{},
				}
			},
			wantErr: true,
			args: args{
				ctx: context.Background(),
			},
		},
		{
			name: "invalid unary expression",
			b: func() BinaryExpression[*MockComparison] {
				return BinaryExpression[*MockComparison]{
					UnaryExpression: &UnaryExpression[*MockComparison]{},
				}
			},
			wantErr: true,
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bexp := tt.b()

			if err := bexp.Validate(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBinaryExpression_Accept(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx          context.Context
		visitorSetup func(visitor *MockVisitor, bexpr *BinaryExpression[*MockComparison])
	}
	type testCase[T Comparison] struct {
		name    string
		b       func() BinaryExpression[T]
		args    args
		wantErr bool
	}
	tests := []testCase[*MockComparison]{
		{
			name: "visit literal expression",
			b: func() BinaryExpression[*MockComparison] {
				return BinaryExpression[*MockComparison]{
					Literal: utils.Ptr(true),
				}
			},
			args: args{
				ctx: context.Background(),
				visitorSetup: func(visitor *MockVisitor, bexpr *BinaryExpression[*MockComparison]) {
					enterBeCall := visitor.EXPECT().EnterExpressionElement(mock.Anything, mock.AnythingOfType("*base.BinaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return(nil).Once()
					literalCall := visitor.EXPECT().Literal(mock.Anything, true).Return(nil).Once().NotBefore(enterBeCall)
					visitor.EXPECT().LeaveExpressionElement(mock.Anything, mock.AnythingOfType("*base.BinaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return().NotBefore(enterBeCall, literalCall)
				},
			},
			wantErr: false,
		},
		{
			name: "visit comparison expression",
			b: func() BinaryExpression[*MockComparison] {
				compMock := NewMockComparison(t)

				return BinaryExpression[*MockComparison]{
					Comparison: compMock,
				}
			},
			args: args{
				ctx: context.Background(),
				visitorSetup: func(visitor *MockVisitor, bexpr *BinaryExpression[*MockComparison]) {
					enterBeCall := visitor.EXPECT().EnterExpressionElement(mock.Anything, mock.AnythingOfType("*base.BinaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return(nil).Once()

					compCall := bexpr.Comparison.EXPECT().Accept(mock.Anything, mock.Anything).Return(nil).Once().NotBefore(enterBeCall)

					visitor.EXPECT().LeaveExpressionElement(mock.Anything, mock.AnythingOfType("*base.BinaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return().NotBefore(compCall)
				},
			},
			wantErr: false,
		},
		{
			name: "visit aggregator expression",
			b: func() BinaryExpression[*MockComparison] {
				return BinaryExpression[*MockComparison]{
					Aggregator: &Aggregator[*MockComparison]{
						Operator: AggregatorOperatorAnd,
						Operands: []BinaryExpression[*MockComparison]{
							{
								Literal: utils.Ptr(true),
							},
						},
					},
				}
			},
			args: args{
				ctx: context.Background(),
				visitorSetup: func(visitor *MockVisitor, bexpr *BinaryExpression[*MockComparison]) {
					visitor.EXPECT().EnterExpressionElement(mock.Anything, mock.AnythingOfType("*base.BinaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return(nil).Twice()
					visitor.EXPECT().EnterExpressionElement(mock.Anything, mock.AnythingOfType("*base.Aggregator[*github.com/collibra/go-bexpression/base.MockComparison]")).Return(nil).Once()
					visitor.EXPECT().Literal(mock.Anything, true).Return(nil).Once()
					visitor.EXPECT().LeaveExpressionElement(mock.Anything, mock.AnythingOfType("*base.Aggregator[*github.com/collibra/go-bexpression/base.MockComparison]")).Return().Once()
					visitor.EXPECT().LeaveExpressionElement(mock.Anything, mock.AnythingOfType("*base.BinaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return().Twice()
				},
			},
		},
		{
			name: "visit unary expression",
			b: func() BinaryExpression[*MockComparison] {
				return BinaryExpression[*MockComparison]{
					UnaryExpression: &UnaryExpression[*MockComparison]{
						Operator: UnaryOperatorNot,
						Operand: BinaryExpression[*MockComparison]{
							Literal: utils.Ptr(true),
						},
					},
				}
			},
			args: args{
				ctx: context.Background(),
				visitorSetup: func(visitor *MockVisitor, bexpr *BinaryExpression[*MockComparison]) {
					visitor.EXPECT().EnterExpressionElement(mock.Anything, mock.AnythingOfType("*base.BinaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return(nil).Twice()
					visitor.EXPECT().EnterExpressionElement(mock.Anything, mock.AnythingOfType("*base.UnaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return(nil).Once()
					visitor.EXPECT().Literal(mock.Anything, UnaryOperatorNot).Return(nil)
					visitor.EXPECT().Literal(mock.Anything, true).Return(nil).Once()
					visitor.EXPECT().LeaveExpressionElement(mock.Anything, mock.AnythingOfType("*base.UnaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return().Once()
					visitor.EXPECT().LeaveExpressionElement(mock.Anything, mock.AnythingOfType("*base.BinaryExpression[*github.com/collibra/go-bexpression/base.MockComparison]")).Return().Twice()
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bexp := tt.b()
			mockVisitor := NewMockVisitor(t)
			tt.args.visitorSetup(mockVisitor, &bexp)

			if err := bexp.Accept(tt.args.ctx, mockVisitor); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBinaryExpression_ToGql(t *testing.T) {
	t.Parallel()

	type testCase[T Comparison] struct {
		name    string
		b       func() BinaryExpression[T]
		want    BinaryExpressionUnion
		wantErr bool
	}
	tests := []testCase[*MockComparison]{
		{
			name: "literal gql",
			b: func() BinaryExpression[*MockComparison] {
				return BinaryExpression[*MockComparison]{
					Literal: utils.Ptr(true),
				}
			},
			want:    &LiteralBool{Value: true},
			wantErr: false,
		},
		{
			name: "comparison gql",
			b: func() BinaryExpression[*MockComparison] {
				compMock := NewMockComparison(t)

				compMock.EXPECT().ToGql().Return(&LiteralBool{Value: true}, nil).Once()

				return BinaryExpression[*MockComparison]{
					Comparison: compMock,
				}
			},
			want:    &LiteralBool{Value: true},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.b()

			got, err := b.ToGql()
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
