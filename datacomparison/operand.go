package datacomparison

import (
	"context"
	"fmt"

	"github.com/collibra/go-bexpression/base"
	"github.com/collibra/go-bexpression/utils"
)

var OperandUnionError = base.NewUnionExpressionError("data_comparison.Operand", []string{"Reference", "Literal"})

type Operand struct {
	Reference *Reference `json:"reference,omitempty" yaml:"reference,omitempty" gqlgen:"reference"`
	Literal   *Literal   `json:"literal,omitempty" yaml:"literal,omitempty" gqlgen:"literal"`
}

func (o *Operand) Validate(ctx context.Context) error {
	if utils.CountNonNil(o.Reference, o.Literal) != 1 {
		return OperandUnionError
	}

	if o.Reference != nil {
		if err := o.Reference.Validate(base.CtxExtendPathAndSetElement(ctx, "reference", o)); err != nil {
			return err
		}
	} else if o.Literal != nil {
		if err := o.Literal.Validate(base.CtxExtendPathAndSetElement(ctx, "literal", o)); err != nil {
			return err
		}
	}

	return nil
}

func (o *Operand) Accept(ctx context.Context, visitor base.Visitor) error {
	err := visitor.EnterExpressionElement(ctx, o)
	if err != nil {
		return fmt.Errorf("entery operand: %w", err)
	}

	defer visitor.LeaveExpressionElement(ctx, o)

	if o.Reference != nil {
		err = o.Reference.Accept(base.CtxExtendPathAndSetElement(ctx, "reference", o), visitor)
		if err != nil {
			return fmt.Errorf("reference: %w", err)
		}

		return nil
	} else if o.Literal != nil {
		err = o.Literal.Accept(base.CtxExtendPathAndSetElement(ctx, "literal", o), visitor)
		if err != nil {
			return fmt.Errorf("literal: %w", err)
		}

		return nil
	}

	return OperandUnionError
}

func (o *Operand) ToGql() (DataComparisonOperand, error) {
	if o.Reference != nil {
		return o.Reference, nil
	} else if o.Literal != nil {
		return o.Literal.ToGql()
	}

	return nil, OperandUnionError
}
