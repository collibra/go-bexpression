package base

import "context"

type Validatable interface {
	Validate(ctx context.Context) error
}

type VisitableElement interface {
	Accept(ctx context.Context, visitor Visitor) error
}

//go:generate go run github.com/vektra/mockery/v3 --name=Visitor --testonly=False
type Visitor interface {
	EnterExpressionElement(ctx context.Context, element VisitableElement) error
	LeaveExpressionElement(ctx context.Context, element VisitableElement)

	Literal(ctx context.Context, l interface{}) error
}

//go:generate go run github.com/vektra/mockery/v3 --name=Comparison
type Comparison interface {
	Validatable
	VisitableElement
	ToGql() (BinaryExpressionUnion, error)
}

type BinaryExpressionUnion interface {
	IsBinaryExpression()
}
