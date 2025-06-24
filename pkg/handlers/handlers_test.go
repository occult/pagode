package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakeInertia struct {
	CalledBack bool
}

func (f *FakeInertia) Back(w http.ResponseWriter, r *http.Request, _ ...int) {
	f.CalledBack = true
}

func TestGetSetHandlers(t *testing.T) {
	handlers = []Handler{}
	assert.Empty(t, GetHandlers())
	h := new(Pages)
	Register(h)
	got := GetHandlers()
	require.Len(t, got, 1)
	assert.Equal(t, h, got[0])
}

func TestFail(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	fake := &FakeInertia{}

	err := fail(errors.New("fail!"), "something went wrong", fake, ctx)

	assert.NoError(t, err)
	assert.True(t, fake.CalledBack, "expected inertia.Back to be called")
}
