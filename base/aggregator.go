package base

import (
	"context"
	"fmt"
)

//go:generate go run github.com/dmarkham/enumer -type=AggregatorOperator -values -gqlgen -yaml -json -trimprefix=AggregatorOperator
type AggregatorOperator int

const (
	AggregatorOperatorAnd AggregatorOperator = iota
	AggregatorOperatorOr
)

type Aggregator[T Comparison] struct {
	Operator AggregatorOperator    `json:"operator" yaml:"operator" gqlgen:"operator"`
	Operands []BinaryExpression[T] `json:"operands" yaml:"operands" gqlgen:"operands"`
}

func (a *Aggregator[T]) Validate(ctx context.Context) error {
	if len(a.Operands) == 0 {
		return ErrEmptyOperands
	}

	for i, operand := range a.Operands {
		if err := operand.Validate(CtxExtendPathAndSetElement(ctx, fmt.Sprintf("operand[%d]", i), a)); err != nil {
			return fmt.Errorf("aggregator[%d]: %w", i, err)
		}
	}

	return nil
}

func (a *Aggregator[T]) Accept(ctx context.Context, visitor Visitor) error {
	err := visitor.EnterExpressionElement(ctx, a)
	if err != nil {
		return fmt.Errorf("enter aggregator: %w", err)
	}

	defer visitor.LeaveExpressionElement(ctx, a)

	for i, operand := range a.Operands {
		if i > 0 {
			err = visitor.Literal(CtxExtendPathAndSetElement(ctx, "operator", a), a.Operator)
			if err != nil {
				return fmt.Errorf("aggregator[%d]: %w", i, err)
			}
		}

		err = operand.Accept(CtxExtendPathAndSetElement(ctx, fmt.Sprintf("operand[%d]", i), a), visitor)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *Aggregator[T]) IsBinaryExpression() {}
