package datacomparison

import (
	"context"
	"fmt"

	"github.com/raito-io/bexpression/base"
)

//go:generate go run github.com/dmarkham/enumer -type=ComparisonOperator -values -gqlgen -yaml -json -trimprefix=ComparisonOperator
type ComparisonOperator int

const (
	ComparisonOperatorEqual ComparisonOperator = iota
	ComparisonOperatorNotEqual
	ComparisonOperatorLessThan
	ComparisonOperatorLessThanOrEqual
	ComparisonOperatorGreaterThan
	ComparisonOperatorGreaterThanOrEqual
)

type DataComparison struct {
	Operator     ComparisonOperator `json:"operator" yaml:"operator" gqlgen:"operator"`
	LeftOperand  Operand            `json:"leftOperand,omitempty" yaml:"leftOperand,omitempty" gqlgen:"leftOperand"`
	RightOperand Operand            `json:"rightOperand,omitempty" yaml:"rightOperand,omitempty" gqlgen:"rightOperand"`
}

func (d *DataComparison) Validate(ctx context.Context) error {
	err := d.LeftOperand.Validate(base.CtxExtendPathAndSetElement(ctx, "leftOperand", d))
	if err != nil {
		return err
	}

	err = d.RightOperand.Validate(base.CtxExtendPathAndSetElement(ctx, "rightOperand", d))
	if err != nil {
		return err
	}

	return nil
}

func (d *DataComparison) Accept(ctx context.Context, visitor base.Visitor) error {
	err := visitor.EnterExpressionElement(ctx, d)
	if err != nil {
		return fmt.Errorf("enter data comparison: %w", err)
	}

	defer visitor.LeaveExpressionElement(ctx, d)

	err = d.LeftOperand.Accept(base.CtxExtendPathAndSetElement(ctx, "leftOperand", d), visitor)
	if err != nil {
		return fmt.Errorf("left operand: %w", err)
	}

	err = visitor.Literal(base.CtxExtendPathAndSetElement(ctx, "operator", d), d.Operator)
	if err != nil {
		return fmt.Errorf("operator: %w", err)
	}

	err = d.RightOperand.Accept(base.CtxExtendPathAndSetElement(ctx, "rightOperand", d), visitor)
	if err != nil {
		return fmt.Errorf("right operand: %w", err)
	}

	return nil
}

func (d *DataComparison) ToGql() (base.BinaryExpressionUnion, error) {
	return d, nil
}

func (d *DataComparison) IsBinaryExpression() {}
