package datacomparison

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/collibra/go-bexpression/base"
)

func TestReference_Accept(t *testing.T) {
	type fields struct {
		EntityType EntityType
		EntityID   string
	}
	type args struct {
		ctx          context.Context
		visitorSetup func(visitor *base.MockVisitor)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				EntityType: EntityTypeDataObject,
				EntityID:   "someId1",
			},
			args: args{
				ctx: context.Background(),
				visitorSetup: func(visitor *base.MockVisitor) {
					visitor.EXPECT().Literal(mock.Anything, &Reference{EntityType: EntityTypeDataObject, EntityID: "someId1"}).Return(nil).Once()
				},
			},
			wantErr: false,
		},
		{
			name: "visitor failed",
			fields: fields{
				EntityType: EntityTypeDataObject,
				EntityID:   "someId1",
			},
			args: args{
				ctx: context.Background(),
				visitorSetup: func(visitor *base.MockVisitor) {
					visitor.EXPECT().Literal(mock.Anything, &Reference{EntityType: EntityTypeDataObject, EntityID: "someId1"}).Return(errors.New("boom")).Once()
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visitorMock := base.NewMockVisitor(t)
			tt.args.visitorSetup(visitorMock)

			r := &Reference{
				EntityType: tt.fields.EntityType,
				EntityID:   tt.fields.EntityID,
			}
			if err := r.Accept(tt.args.ctx, visitorMock); (err != nil) != tt.wantErr {
				t.Errorf("Accept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReference_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ref     Reference
		wantErr error
	}{
		{name: "empty reference", ref: Reference{}},
		{name: "column reference by name", ref: Reference{EntityType: EntityTypeColumnReferenceByName, EntityID: "COUNTRY"}},
		{name: "tag reference with key", ref: Reference{EntityType: EntityTypeColumnReferenceByTag, TagKey: "Country"}},
		{name: "tag reference with key and value", ref: Reference{EntityType: EntityTypeColumnReferenceByTag, TagKey: "Country", TagValue: "DE"}},
		{name: "tag reference without key", ref: Reference{EntityType: EntityTypeColumnReferenceByTag}, wantErr: ErrTagKeyRequired},
		{name: "tag reference with value only", ref: Reference{EntityType: EntityTypeColumnReferenceByTag, TagValue: "DE"}, wantErr: ErrTagKeyRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ref.Validate(context.Background())

			if tt.wantErr == nil {
				assert.NoError(t, err)
				return
			}
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestReference_JSONRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		ref      Reference
		wantJSON string
	}{
		{
			name:     "tag key and value",
			ref:      Reference{EntityType: EntityTypeColumnReferenceByTag, TagKey: "Country", TagValue: "DE"},
			wantJSON: `{"entityType":"ColumnReferenceByTag","tagKey":"Country","tagValue":"DE"}`,
		},
		{
			name:     "tag key only omits the value",
			ref:      Reference{EntityType: EntityTypeColumnReferenceByTag, TagKey: "Country"},
			wantJSON: `{"entityType":"ColumnReferenceByTag","tagKey":"Country"}`,
		},
		{
			name:     "column reference by name omits the tag fields",
			ref:      Reference{EntityType: EntityTypeColumnReferenceByName, EntityID: "COUNTRY"},
			wantJSON: `{"entityType":"ColumnReferenceByName","entityId":"COUNTRY"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.ref)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(got))

			var back Reference
			assert.NoError(t, json.Unmarshal(got, &back))
			assert.Equal(t, tt.ref, back)
		})
	}
}

func TestEntityType_ColumnReferenceByTag(t *testing.T) {
	assert.Equal(t, "ColumnReferenceByTag", EntityTypeColumnReferenceByTag.String())

	got, err := EntityTypeString("ColumnReferenceByTag")
	assert.NoError(t, err)
	assert.Equal(t, EntityTypeColumnReferenceByTag, got)
}

func TestReference_Accept_TagReference(t *testing.T) {
	ref := &Reference{EntityType: EntityTypeColumnReferenceByTag, TagKey: "Country", TagValue: "DE"}
	visitorMock := base.NewMockVisitor(t)
	visitorMock.EXPECT().Literal(mock.Anything, ref).Return(nil).Once()

	assert.NoError(t, ref.Accept(context.Background(), visitorMock))
}
