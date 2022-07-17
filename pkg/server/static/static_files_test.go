package static

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNormalize(t *testing.T) {

	type test struct {
		desc       string
		path       string
		normalized string
	}

	for _, tc := range []test{
		{
			desc:       "empty",
			path:       "",
			normalized: "",
		},
		{
			desc:       "regular path",
			path:       "/test/index.html",
			normalized: "test/index.html",
		},
		{
			desc:       "path with slashes",
			path:       "\\test\\index.html",
			normalized: "test/index.html",
		},
		{
			desc:       "path with parent directory",
			path:       "/test/../index.html",
			normalized: "test/index.html",
		},
		{
			desc:       "path with current directory",
			path:       "/./test/index.html",
			normalized: "test/index.html",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			normalized := normalize(tc.path)
			require.Equal(t, tc.normalized, normalized)
		})
	}
}
