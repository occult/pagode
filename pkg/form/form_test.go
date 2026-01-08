package form

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/pkg/context"
	"github.com/occult/pagode/pkg/tests"
	"github.com/romsar/gonertia/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockForm struct {
	called bool
	Submission
}

func (m *mockForm) Submit(_ echo.Context, _ any) error {
	m.called = true
	return nil
}

func TestSubmit(t *testing.T) {
	m := mockForm{}
	ctx, _ := tests.NewContext(echo.New(), "/")
	err := Submit(ctx, &m)
	require.NoError(t, err)
	assert.True(t, m.called)
}

func TestGetClear(t *testing.T) {
	e := echo.New()

	type example struct {
		Name string `form:"name"`
	}

	t.Run("get empty context", func(t *testing.T) {
		// Empty context, still return a form
		ctx, _ := tests.NewContext(e, "/")
		form := Get[example](ctx)
		assert.NotNil(t, form)
	})

	t.Run("get non-empty context", func(t *testing.T) {
		form := example{
			Name: "test",
		}
		ctx, _ := tests.NewContext(e, "/")
		ctx.Set(context.FormKey, &form)

		// Get again and expect the values were stored
		got := Get[example](ctx)
		require.NotNil(t, got)
		assert.Equal(t, "test", form.Name)

		// Clear
		Clear(ctx)
		got = Get[example](ctx)
		require.NotNil(t, got)
		assert.Empty(t, got.Name)
	})
}

func TestShareErrorsPreservesExistingProps(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("POST", "/test", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Set up existing Inertia props
	existingProps := map[string]any{
		"auth": map[string]any{
			"user": map[string]any{"name": "Test User"},
		},
		"flash": map[string]any{
			"success": []string{"Welcome!"},
		},
	}
	newReqCtx := gonertia.SetProps(req.Context(), existingProps)
	ctx.SetRequest(req.WithContext(newReqCtx))

	// Create a form with errors
	form := &mockForm{}
	form.SetFieldError("Email", "Email is required")
	form.SetFieldError("Password", "Password is required")

	// Share errors
	ShareErrors(ctx, form)

	// Get props from context
	props := gonertia.PropsFromContext(ctx.Request().Context())
	require.NotNil(t, props)

	// Verify existing props are preserved
	auth, ok := props["auth"].(map[string]any)
	require.True(t, ok, "auth prop should exist")
	user, ok := auth["user"].(map[string]any)
	require.True(t, ok, "auth.user should exist")
	assert.Equal(t, "Test User", user["name"])

	flash, ok := props["flash"].(map[string]any)
	require.True(t, ok, "flash prop should exist")
	success, ok := flash["success"].([]string)
	require.True(t, ok, "flash.success should exist")
	assert.Equal(t, []string{"Welcome!"}, success)

	// Verify errors are added
	errors, ok := props["errors"].(map[string][]string)
	require.True(t, ok, "errors prop should exist")
	assert.Equal(t, []string{"Email is required"}, errors["Email"])
	assert.Equal(t, []string{"Password is required"}, errors["Password"])
}

func TestShareErrorsWorksWithNoExistingProps(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("POST", "/test", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Create a form with errors
	form := &mockForm{}
	form.SetFieldError("Name", "Name is required")

	// Share errors (no existing props)
	ShareErrors(ctx, form)

	// Get props from context
	props := gonertia.PropsFromContext(ctx.Request().Context())
	require.NotNil(t, props)

	// Verify errors are added
	errors, ok := props["errors"].(map[string][]string)
	require.True(t, ok, "errors prop should exist")
	assert.Equal(t, []string{"Name is required"}, errors["Name"])
}

func TestShareErrorsNoOpWhenNoErrors(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("POST", "/test", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Set up existing Inertia props
	existingProps := map[string]any{
		"auth": map[string]any{
			"user": "test",
		},
	}
	newReqCtx := gonertia.SetProps(req.Context(), existingProps)
	ctx.SetRequest(req.WithContext(newReqCtx))

	// Create a form with NO errors
	form := &mockForm{}

	// Share errors (should be no-op)
	ShareErrors(ctx, form)

	// Get props from context
	props := gonertia.PropsFromContext(ctx.Request().Context())
	require.NotNil(t, props)

	// Verify existing props are untouched
	assert.Equal(t, "test", props["auth"].(map[string]any)["user"])

	// Verify no errors key was added
	_, hasErrors := props["errors"]
	assert.False(t, hasErrors, "errors prop should not exist when there are no errors")
}
