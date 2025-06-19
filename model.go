package bexpression

import (
	"github.com/collibra/go-bexpression/base"
	"github.com/collibra/go-bexpression/datacomparison"
)

type DataComparisonExpression = base.BinaryExpression[*datacomparison.DataComparison]
type DataComparisonAggregator = base.Aggregator[*datacomparison.DataComparison]
type DataComparisonUnaryExpression = base.UnaryExpression[*datacomparison.DataComparison]
