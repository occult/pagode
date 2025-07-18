package htmx

import (
	"net/http"
	"testing"

	"github.com/occult/pagode/pkg/context"
	"github.com/occult/pagode/pkg/tests"

	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo/v4"
)

func TestSetRequest(t *testing.T) {
	ctx, _ := tests.NewContext(echo.New(), "/")
	ctx.Request().Header.Set(HeaderRequest, "true")
	ctx.Request().Header.Set(HeaderBoosted, "true")
	ctx.Request().Header.Set(HeaderTrigger, "a")
	ctx.Request().Header.Set(HeaderTriggerName, "b")
	ctx.Request().Header.Set(HeaderTarget, "c")
	ctx.Request().Header.Set(HeaderPrompt, "d")
	ctx.Request().Header.Set(HeaderHistoryRestoreRequest, "true")

	r := GetRequest(ctx)
	assert.Equal(t, true, r.Enabled)
	assert.Equal(t, true, r.Boosted)
	assert.Equal(t, true, r.HistoryRestore)
	assert.Equal(t, "a", r.Trigger)
	assert.Equal(t, "b", r.TriggerName)
	assert.Equal(t, "c", r.Target)
	assert.Equal(t, "d", r.Prompt)

	cached := context.Cache(ctx, context.HTMXRequestKey, func(ctx echo.Context) *Request {
		return nil
	})
	assert.Equal(t, r, cached)
}

func TestResponse_Apply(t *testing.T) {
	ctx, _ := tests.NewContext(echo.New(), "/")
	r := Response{
		PushURL:            "a",
		Redirect:           "b",
		ReplaceURL:         "f",
		Refresh:            true,
		Trigger:            "c",
		TriggerAfterSwap:   "d",
		TriggerAfterSettle: "e",
		NoContent:          true,
	}
	r.Apply(ctx)

	assert.Equal(t, "a", ctx.Response().Header().Get(HeaderPushURL))
	assert.Equal(t, "b", ctx.Response().Header().Get(HeaderRedirect))
	assert.Equal(t, "true", ctx.Response().Header().Get(HeaderRefresh))
	assert.Equal(t, "c", ctx.Response().Header().Get(HeaderTrigger))
	assert.Equal(t, "d", ctx.Response().Header().Get(HeaderTriggerAfterSwap))
	assert.Equal(t, "e", ctx.Response().Header().Get(HeaderTriggerAfterSettle))
	assert.Equal(t, "f", ctx.Response().Header().Get(HeaderReplaceURL))
	assert.Equal(t, http.StatusNoContent, ctx.Response().Status)
}
