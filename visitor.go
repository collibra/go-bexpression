package bexpression

import (
	"context"

	"github.com/collibra/go-bexpression/base"
)

type FunctionVisitorOptions struct {
	EnterExpressionElementFn func(ctx context.Context, element base.VisitableElement) error
	LeaveExpressionElementFn func(ctx context.Context, element base.VisitableElement)
	LiteralFn                func(ctx context.Context, literal interface{}) error
}

func WithEnterExpressionElementFn(enterExpressionElementFn func(ctx context.Context, element base.VisitableElement) error) func(*FunctionVisitorOptions) {
	return func(options *FunctionVisitorOptions) {
		options.EnterExpressionElementFn = enterExpressionElementFn
	}
}

func WithLeaveExpressionElementFn(leaveExpressionElementFn func(ctx context.Context, element base.VisitableElement)) func(options *FunctionVisitorOptions) {
	return func(options *FunctionVisitorOptions) {
		options.LeaveExpressionElementFn = leaveExpressionElementFn
	}
}

func WithLiteralFn(literalFn func(ctx context.Context, literal interface{}) error) func(*FunctionVisitorOptions) {
	return func(options *FunctionVisitorOptions) {
		options.LiteralFn = literalFn
	}
}

var _ base.Visitor = (*FunctionVisitor)(nil)

type FunctionVisitor struct {
	EnterExpressionElementFn func(ctx context.Context, element base.VisitableElement) error
	LeaveExpressionElementFn func(ctx context.Context, element base.VisitableElement)
	LiteralFn                func(ctx context.Context, literal interface{}) error
}

func NewFunctionVisitor(opts ...func(*FunctionVisitorOptions)) *FunctionVisitor {
	options := &FunctionVisitorOptions{
		EnterExpressionElementFn: func(ctx context.Context, element base.VisitableElement) error { return nil },
		LeaveExpressionElementFn: func(ctx context.Context, element base.VisitableElement) {},
		LiteralFn:                func(ctx context.Context, literal interface{}) error { return nil },
	}

	for _, opt := range opts {
		opt(options)
	}

	return &FunctionVisitor{
		EnterExpressionElementFn: options.EnterExpressionElementFn,
		LeaveExpressionElementFn: options.LeaveExpressionElementFn,
		LiteralFn:                options.LiteralFn,
	}
}

func (f FunctionVisitor) EnterExpressionElement(ctx context.Context, element base.VisitableElement) error {
	return f.EnterExpressionElementFn(ctx, element)
}

func (f FunctionVisitor) LeaveExpressionElement(ctx context.Context, element base.VisitableElement) {
	f.LeaveExpressionElementFn(ctx, element)
}

func (f FunctionVisitor) Literal(ctx context.Context, l interface{}) error {
	return f.LiteralFn(ctx, l)
}
