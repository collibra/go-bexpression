package datacomparison

import (
	"context"
	"errors"
	"fmt"

	"github.com/collibra/go-bexpression/base"
)

//go:generate go run github.com/dmarkham/enumer -type=EntityType -values -gqlgen -yaml -json -trimprefix=EntityType
type EntityType int

const (
	EntityTypeDataObject EntityType = iota
	EntityTypeColumnReferenceByName
	EntityTypeColumnReferenceByTag
)

var ErrTagKeyRequired = errors.New("tagKey is required for a ColumnReferenceByTag reference")

type Reference struct {
	EntityType EntityType `json:"entityType,omitempty" yaml:"entityType,omitempty" gqlgen:"entityType"`
	EntityID   string     `json:"entityId,omitempty" yaml:"entityId,omitempty" gqlgen:"entityId"`
	// TagKey and TagValue select the column for EntityTypeColumnReferenceByTag. An empty TagValue matches any value.
	TagKey   string `json:"tagKey,omitempty" yaml:"tagKey,omitempty" gqlgen:"tagKey"`
	TagValue string `json:"tagValue,omitempty" yaml:"tagValue,omitempty" gqlgen:"tagValue"`
}

func (r *Reference) Validate(_ context.Context) error {
	if r.EntityType == EntityTypeColumnReferenceByTag && r.TagKey == "" {
		return ErrTagKeyRequired
	}

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
