<h1 align="center">
  Binary Expression Parser
</h1>

<p align="center">
</p>

<hr/>

# Introduction
A library to parse binary expression within Collibra Access Governance plugins and services.
The library provides a way to define binary expressions, marhalize them to JSON, and unmarshal them back to Go structs.
The library also contains a visitor pattern to traverse the binary expression tree.

The language definition can be found here [Binary Expression Language](language)

# Usage

```go
package main

import (
	"context"

	"github.com/collibra/go-bexpression"
	"github.com/collibra/go-bexpression/base"
	"github.com/collibra/go-bexpression/datacomparison"
)

func main() {
	expr := bexpression.DataComparisonExpression{
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

	var enterElements []base.VisitableElement
	var leaveElements []base.VisitableElement
	var literals []interface{}

	visitor := bexpression.NewFunctionVisitor(bexpression.WithEnterExpressionElementFn(
		func(ctx context.Context, element base.VisitableElement) error {
			enterElements = append(enterElements, element)

			return nil
		}), bexpression.WithLeaveExpressionElementFn(func(ctx context.Context, element base.VisitableElement) {
		leaveElements = append(leaveElements, element)
	}), bexpression.WithLiteralFn(func(ctx context.Context, literal interface{}) error {
		literals = append(literals, literal)

		return nil
	}))

	err := expr.Accept(context.Background(), visitor)
	if err != nil {
		panic(err)
    }
	
	// enterElements contains [&expr, expr.Comparison, &expr.Comparison.LeftOperand, &expr.Comparison.RightOperand]
	// leaveElements contains [&expr.Comparison.LeftOperand, &expr.Comparison.RightOperand, expr.Comparison, &expr]
	// literals contains expr.Comparison.LeftOperand.Reference, datacomparison.ComparisonOperatorGreaterThan, float64(3.14),
}
```