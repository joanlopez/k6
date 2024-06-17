//go:build wpt

package streams

import (
	"testing"

	"go.k6.io/k6/js/modules/k6/timers"
	"go.k6.io/k6/js/modulestest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlob(t *testing.T) {
	t.Parallel()

	suites := []string{
		"Blob-text.any.js",
		"Blob-array-buffer.any.js",
	}

	for _, suite := range suites {
		suite := suite
		t.Run(suite, func(t *testing.T) {
			t.Parallel()
			ts := newConfiguredRuntimeForBlob(t)
			gotErr := ts.EventLoop.Start(func() error {
				return executeTestScript(ts.VU, "tests/wpt/FileAPI/blob", suite)
			})
			assert.NoError(t, gotErr)
		})
	}
}

func newConfiguredRuntimeForBlob(t testing.TB) *modulestest.Runtime {
	rt := modulestest.NewRuntime(t)

	// We want to make the [self] available for Web Platform Tests, as it is used in test harness.
	_, err := rt.VU.Runtime().RunString("var self = this;")
	require.NoError(t, err)

	// We want to make the [console.log()] available for Web Platform Tests, as it
	// is very useful for debugging, because we don't have a real debugger for JS code.
	logger := rt.VU.InitEnvField.Logger
	require.NoError(t, rt.VU.RuntimeField.Set("console", newConsole(logger)))

	// We also want to make [timers.Timers] available for Web Platform Tests.
	for k, v := range timers.New().NewModuleInstance(rt.VU).Exports().Named {
		require.NoError(t, rt.VU.RuntimeField.Set(k, v))
	}

	// We also want the Blob module exports to be globally available.
	m := new(RootModuleBlob).NewModuleInstance(rt.VU)
	for k, v := range m.Exports().Named {
		require.NoError(t, rt.VU.RuntimeField.Set(k, v))
	}

	// Then, we register the Web Platform Tests harness.
	compileAndRun(t, rt, "tests/wpt", "resources/testharness.js")

	// And the Blob-specific test utilities.
	files := []string{
		"support/Blob.js",
	}
	for _, file := range files {
		compileAndRun(t, rt, "tests/wpt/FileAPI", file)
	}

	return rt
}
