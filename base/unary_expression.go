package base

import (
	"context"
	"fmt"
)

//go:generate go run github.com/dmarkham/enumer -type=UnaryOperator -values -gqlgen -yaml -json -trimprefix=UnaryOperator
type UnaryOperator int

const (
	UnaryOperatorNot UnaryOperator = iota
)

type UnaryExpression[T Comparison] struct {
	Operator UnaryOperator       `json:"operator" yaml:"operator" gqlgen:"operator"`
	Operand  BinaryExpression[T] `json:"expression" yaml:"expression" gqlgen:"expression"`
}

func (n *UnaryExpression[T]) Validate(ctx context.Context) error {
	return n.Operand.Validate(CtxExtendPathAndSetElement(ctx, "operand", n))
}

func (n *UnaryExpression[T]) Accept(ctx context.Context, visitor Visitor) error {
	err := visitor.EnterExpressionElement(ctx, n)
	if err != nil {
		return fmt.Errorf("enter unary expression: %w", err)
	}

	defer visitor.LeaveExpressionElement(ctx, n)

	err = visitor.Literal(CtxExtendPathAndSetElement(ctx, "operator", n), n.Operator)
	if err != nil {
		return fmt.Errorf("operator: %w", err)
	}

	err = n.Operand.Accept(CtxExtendPathAndSetElement(ctx, "operand", n), visitor)
	if err != nil {
		return fmt.Errorf("operand: %w", err)
	}

	return nil
}

func (n *UnaryExpression[T]) IsBinaryExpression() {}
