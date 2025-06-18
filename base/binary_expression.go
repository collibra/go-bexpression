package base

import (
	"context"
	"errors"
	"fmt"

	"github.com/collibra/go-bexpression/utils"
)

var (
	BinaryExpressionUnionError = NewUnionExpressionError("BinaryExpression", []string{"Literal", "Comparison", "Aggregator", "UnaryExpression"})
)

type BinaryExpression[T Comparison] struct {
	Literal         *bool               `json:"literal,omitempty" yaml:"literal,omitempty" gqlgen:"literal"`
	Comparison      T                   `json:"comparison,omitempty" yaml:"comparison,omitempty" gqlgen:"comparison"`
	Aggregator      *Aggregator[T]      `json:"aggregator,omitempty" yaml:"aggregator,omitempty" gqlgen:"aggregator"`
	UnaryExpression *UnaryExpression[T] `json:"unaryExpression,omitempty" yaml:"unaryExpression,omitempty" gqlgen:"unaryExpression"`
}

func (b *BinaryExpression[T]) Validate(ctx context.Context) error {
	if utils.CountNonNil(b.Literal, b.Comparison, b.Aggregator, b.UnaryExpression) != 1 {
		return BinaryExpressionUnionError
	}

	if b.Literal != nil {
		return nil
	} else if !utils.IsNull(b.Comparison) {
		err := b.Comparison.Validate(CtxExtendPathAndSetElement(ctx, "comparison", b))
		if err != nil {
			return fmt.Errorf("comparison: %w", err)
		}

		return nil
	} else if b.Aggregator != nil {
		err := b.Aggregator.Validate(CtxExtendPathAndSetElement(ctx, "aggregator", b))
		if err != nil {
			return fmt.Errorf("aggregator: %w", err)
		}

		return nil
	} else if b.UnaryExpression != nil {
		err := b.UnaryExpression.Validate(CtxExtendPathAndSetElement(ctx, "unaryExpression", b))
		if err != nil {
			return fmt.Errorf("unaryExpression: %w", err)
		}

		return nil
	}

	return errors.New("no valid expression element set") // This line may never be executed
}

func (b *BinaryExpression[T]) Accept(ctx context.Context, visitor Visitor) error {
	err := visitor.EnterExpressionElement(ctx, b)
	if err != nil {
		return fmt.Errorf("enter binary expression: %w", err)
	}

	defer visitor.LeaveExpressionElement(ctx, b)

	if b.Literal != nil {
		err = visitor.Literal(CtxExtendPathAndSetElement(ctx, "literal", b), *b.Literal)
		if err != nil {
			return fmt.Errorf("literal: %w", err)
		}

		return nil
	} else if !utils.IsNull(b.Comparison) {
		err = b.Comparison.Accept(CtxExtendPathAndSetElement(ctx, "comparison", b), visitor)
		if err != nil {
			return fmt.Errorf("comparison: %w", err)
		}

		return nil
	} else if b.Aggregator != nil {
		err = b.Aggregator.Accept(CtxExtendPathAndSetElement(ctx, "aggregator", b), visitor)
		if err != nil {
			return fmt.Errorf("aggregator: %w", err)
		}

		return nil
	} else if b.UnaryExpression != nil {
		err = b.UnaryExpression.Accept(CtxExtendPathAndSetElement(ctx, "unaryExpression", b), visitor)
		if err != nil {
			return fmt.Errorf("unaryExpression: %w", err)
		}

		return nil
	}

	return BinaryExpressionUnionError
}

func (b *BinaryExpression[T]) ToGql() (BinaryExpressionUnion, error) {
	if b.Literal != nil {
		return &LiteralBool{Value: *b.Literal}, nil
	} else if !utils.IsNull(b.Comparison) {
		r, err := b.Comparison.ToGql()
		if err != nil {
			return nil, fmt.Errorf("comparison: %w", err)
		}

		return r, nil
	} else if b.Aggregator != nil {
		return b.Aggregator, nil
	} else if b.UnaryExpression != nil {
		return b.UnaryExpression, nil
	}

	return nil, BinaryExpressionUnionError
}
