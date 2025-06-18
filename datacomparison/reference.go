package datacomparison

import (
	"context"
	"fmt"

	"github.com/raito-io/bexpression/base"
)

//go:generate go run github.com/dmarkham/enumer -type=EntityType -values -gqlgen -yaml -json -trimprefix=EntityType
type EntityType int

const (
	EntityTypeDataObject EntityType = iota
	EntityTypeColumnReferenceByName
)

type Reference struct {
	EntityType EntityType `json:"entityType,omitempty" yaml:"entityType,omitempty" gqlgen:"entityType"`
	EntityID   string     `json:"entityId,omitempty" yaml:"entityId,omitempty" gqlgen:"entityId"`
}

func (r *Reference) Validate(_ context.Context) error {
	return nil
}

func (r *Reference) Accept(ctx context.Context, visitor base.Visitor) error {
	err := visitor.Literal(ctx, r)
	if err != nil {
		return fmt.Errorf("literal: %w", err)
	}

	return nil
}

func (r *Reference) IsDataComparisonOperand() {}
