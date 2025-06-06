package limit

import (
	"context"
	"github.com/anyproto/any-sync/app"
	"github.com/anyproto/any-sync/coordinator/coordinatorclient"
	"github.com/anyproto/any-sync/coordinator/coordinatorclient/mock_coordinatorclient"
	"github.com/anyproto/any-sync/net/peer"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var ctx = context.Background()

func TestLimit_Check(t *testing.T) {
	var spaceId = "122345.123"
	var identity = []byte("identity")
	t.Run("succes", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish(t)
		fx.client.EXPECT().FileLimitCheck(gomock.Any(), spaceId, identity).Return(uint64(123), nil)
		res, err := fx.Check(peer.CtxWithIdentity(ctx, identity), spaceId)
		require.NoError(t, err)
		assert.Equal(t, uint64(123), res)
		res, err = fx.Check(peer.CtxWithIdentity(ctx, identity), spaceId)
		require.NoError(t, err)
		assert.Equal(t, uint64(123), res)
	})
	t.Run("no identity", func(t *testing.T) {
		fx := newFixture(t)
		defer fx.Finish(t)
		_, err := fx.Check(ctx, spaceId)
		require.Error(t, err)
	})

}

func newFixture(t *testing.T) *fixture {
	ctrl := gomock.NewController(t)
	fx := &fixture{
		Limit:  New(),
		client: mock_coordinatorclient.NewMockCoordinatorClient(ctrl),
		ctrl:   ctrl,
		a:      new(app.App),
	}
	fx.client.EXPECT().Name().Return(coordinatorclient.CName).AnyTimes()
	fx.client.EXPECT().Init(gomock.Any()).AnyTimes()
	fx.a.Register(fx.client).Register(fx.Limit)
	require.NoError(t, fx.a.Start(ctx))
	return fx
}

type fixture struct {
	Limit
	client *mock_coordinatorclient.MockCoordinatorClient
	ctrl   *gomock.Controller
	a      *app.App
}

func (fx *fixture) Finish(t *testing.T) {
	require.NoError(t, fx.a.Close(ctx))
	fx.ctrl.Finish()
}
