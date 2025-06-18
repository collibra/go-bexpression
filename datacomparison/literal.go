package datacomparison

import (
	"context"
	"fmt"
	"time"

	"github.com/collibra/go-bexpression/base"
	"github.com/collibra/go-bexpression/utils"
)

var LiteralUnionErr = base.NewUnionExpressionError("dataComparison.Literal", []string{"bool", "int", "float", "string", "timestamp"})

type Literal struct {
	Bool      *bool      `json:"bool,omitempty" yaml:"bool,omitempty" gqlgen:"bool"`
	Int       *int       `json:"int,omitempty" yaml:"int,omitempty" gqlgen:"int"`
	Float     *float64   `json:"float,omitempty" yaml:"float,omitempty" gqlgen:"float"`
	Str       *string    `json:"string,omitempty" yaml:"string,omitempty" gqlgen:"string"`
	Timestamp *time.Time `json:"timestamp,omitempty" yaml:"timestamp,omitempty" gqlgen:"timestamp"`
}

func (l *Literal) Validate(_ context.Context) error {
	if utils.CountNonNil(l.Bool, l.Int, l.Float, l.Str, l.Timestamp) != 1 {
		return LiteralUnionErr
	}

	return nil
}

func (l *Literal) Accept(ctx context.Context, visitor base.Visitor) error {
	if l.Bool != nil {
		err := visitor.Literal(base.CtxExtendPathAndSetElement(ctx, "bool", l), *l.Bool)
		if err != nil {
			return fmt.Errorf("bool: %w", err)
		}

		return nil
	} else if l.Int != nil {
		err := visitor.Literal(base.CtxExtendPathAndSetElement(ctx, "int", l), *l.Int)
		if err != nil {
			return fmt.Errorf("int: %w", err)
		}

		return nil
	} else if l.Float != nil {
		err := visitor.Literal(base.CtxExtendPathAndSetElement(ctx, "float", l), *l.Float)
		if err != nil {
			return fmt.Errorf("float: %w", err)
		}

		return nil
	} else if l.Str != nil {
		err := visitor.Literal(base.CtxExtendPathAndSetElement(ctx, "string", l), *l.Str)
		if err != nil {
			return fmt.Errorf("string: %w", err)
		}

		return nil
	} else if l.Timestamp != nil {
		err := visitor.Literal(base.CtxExtendPathAndSetElement(ctx, "timestamp", l), *l.Timestamp)
		if err != nil {
			return fmt.Errorf("timestamp: %w", err)
		}

		return nil
	}

	return LiteralUnionErr
}

func (l *Literal) ToGql() (DataComparisonOperand, error) {
	if l.Bool != nil {
		return &LiteralBool{Value: *l.Bool}, nil
	} else if l.Int != nil {
		return &LiteralInt{Value: *l.Int}, nil
	} else if l.Float != nil {
		return &LiteralFloat{Value: *l.Float}, nil
	} else if l.Str != nil {
		return &LiteralString{Value: *l.Str}, nil
	} else if l.Timestamp != nil {
		return &LiteralTime{Value: *l.Timestamp}, nil
	}

	return nil, LiteralUnionErr
}

type LiteralBool base.Literal[bool]

func (l *LiteralBool) IsDataComparisonOperand() {}

type LiteralInt base.Literal[int]

func (l *LiteralInt) IsDataComparisonOperand() {}

type LiteralFloat base.Literal[float64]

func (l *LiteralFloat) IsDataComparisonOperand() {}

type LiteralString base.Literal[string]

func (l *LiteralString) IsDataComparisonOperand() {}

type LiteralTime base.Literal[time.Time]

func (l *LiteralTime) IsDataComparisonOperand() {}
