package bexpression

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/collibra/go-bexpression/base"
	"github.com/collibra/go-bexpression/datacomparison"
	"github.com/collibra/go-bexpression/utils"
)

func TestFunctionVisitor(t *testing.T) {
	expr := DataComparisonExpression{
		Comparison: &datacomparison.DataComparison{
			Operator: datacomparison.ComparisonOperatorGreaterThan,
			LeftOperand: datacomparison.Operand{
				Reference: &datacomparison.Reference{
					EntityType: datacomparison.EntityTypeDataObject,
					EntityID:   "someEntityId",
				},
			},
			RightOperand: datacomparison.Operand{
				Literal: &datacomparison.Literal{
					Float: utils.Ptr(3.14),
				},
			},
		},
	}

	t.Run("Fully initialized visitor", func(t *testing.T) {
		// Given
		var enterElements []base.VisitableElement
		var leaveElements []base.VisitableElement
		var literals []interface{}

		visitor := NewFunctionVisitor(WithEnterExpressionElementFn(func(ctx context.Context, element base.VisitableElement) error {
			enterElements = append(enterElements, element)

			return nil
		}), WithLeaveExpressionElementFn(func(ctx context.Context, element base.VisitableElement) {
			leaveElements = append(leaveElements, element)
		}), WithLiteralFn(func(ctx context.Context, literal interface{}) error {
			literals = append(literals, literal)

			return nil
		}))

		// When
		err := expr.Accept(context.Background(), visitor)

		// Then
		require.NoError(t, err)

		assert.Equal(t, []base.VisitableElement{
			&expr, expr.Comparison, &expr.Comparison.LeftOperand, &expr.Comparison.RightOperand,
		}, enterElements)

		assert.Equal(t, []base.VisitableElement{
			&expr.Comparison.LeftOperand, &expr.Comparison.RightOperand, expr.Comparison, &expr,
		}, leaveElements)

		assert.Equal(t, []interface{}{
			expr.Comparison.LeftOperand.Reference, datacomparison.ComparisonOperatorGreaterThan, float64(3.14),
		}, literals)
	})

	t.Run("Empty visitor", func(t *testing.T) {
		// Given
		visitor := NewFunctionVisitor()

		// When
		err := expr.Accept(context.Background(), visitor)

		// Then
		require.NoError(t, err)
	})
}
