package datacomparison

import (
	"context"
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
	r := Reference{}

	assert.NoError(t, r.Validate(context.Background()))
}
