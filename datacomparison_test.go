package bexpression

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/collibra/go-bexpression/base"
	"github.com/collibra/go-bexpression/datacomparison"
	"github.com/collibra/go-bexpression/utils"
)

var _ (base.Visitor) = (*testVisitor)(nil)

type Stringify interface {
	String() string
}

type testVisitor struct {
	strBuilder strings.Builder

	t string
}

func (t *testVisitor) EnterExpressionElement(_ context.Context, element base.VisitableElement) error {
	switch element.(type) {
	case *base.BinaryExpression[*datacomparison.DataComparison]:
		t.strBuilder.WriteString("(")
	}

	return nil
}

func (t *testVisitor) LeaveExpressionElement(_ context.Context, element base.VisitableElement) {
	switch element.(type) {
	case *base.BinaryExpression[*datacomparison.DataComparison]:
		t.strBuilder.WriteString(")")
	}
}

func (t *testVisitor) Literal(_ context.Context, l interface{}) error {
	if s, ok := l.(Stringify); ok {
		t.strBuilder.WriteString(" " + s.String() + " ")

	} else if reference, ok := l.(*datacomparison.Reference); ok {
		t.strBuilder.WriteString(fmt.Sprintf("%s.%s", reference.EntityType.String(), reference.EntityID))

	} else {
		t.strBuilder.WriteString(fmt.Sprintf("%v ", l))
	}

	return nil
}

func TestDataComparisonExpression_Visitor(t *testing.T) {
	tests := []struct {
		name string
		expr *DataComparisonExpression
		want string
	}{
		{
			name: "simple comparison",
			expr: &DataComparisonExpression{
				Comparison: &datacomparison.DataComparison{
					Operator: datacomparison.ComparisonOperatorGreaterThan,
					LeftOperand: datacomparison.Operand{
						Literal: &datacomparison.Literal{
							Int: utils.Ptr(10),
						},
					},
					RightOperand: datacomparison.Operand{
						Reference: &datacomparison.Reference{
							EntityType: datacomparison.EntityTypeDataObject,
							EntityID:   "EntityID",
						},
					},
				},
			},
			want: "(10  GreaterThan DataObject.EntityID)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visitor := testVisitor{
				strBuilder: strings.Builder{},
				t:          "root",
			}
			err := tt.expr.Accept(context.Background(), &visitor)

			require.NoError(t, err)

			if visitor.strBuilder.String() != tt.want {
				t.Errorf("DataComparisonExpression.Visitor() = %v, want %v", visitor.strBuilder.String(), tt.want)
			}
		})
	}
}
