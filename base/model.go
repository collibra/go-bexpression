package base

import "context"

type Validatable interface {
	Validate(ctx context.Context) error
}

type VisitableElement interface {
	Accept(ctx context.Context, visitor Visitor) error
}

type Visitor interface {
	EnterExpressionElement(ctx context.Context, element VisitableElement) error
	LeaveExpressionElement(ctx context.Context, element VisitableElement)

	Literal(ctx context.Context, l interface{}) error
}

type Comparison interface {
	Validatable
	VisitableElement
	ToGql() (BinaryExpressionUnion, error)
}

type BinaryExpressionUnion interface {
	IsBinaryExpression()
}
