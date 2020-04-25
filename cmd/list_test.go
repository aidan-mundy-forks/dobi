package cmd

import (
	"testing"

	"github.com/dnephin/dobi/config"
	testconfig "github.com/dnephin/dobi/internal/test/config"
	"gotest.tools/v3/assert"
	is "gotest.tools/v3/assert/cmp"
)

func TestInclude(t *testing.T) {
	var testcases = []struct {
		doc      string
		opts     listOptions
		resource config.Resource
		expected bool
	}{
		{
			resource: &testconfig.FakeResource{},
			expected: false,
		},
		{
			opts:     listOptions{all: true},
			resource: &testconfig.FakeResource{},
			expected: true,
		},
		{
			opts:     listOptions{tags: []string{"one"}},
			resource: &testconfig.FakeResource{},
			expected: false,
		},
		{
			opts: listOptions{tags: []string{"one"}},
			resource: &testconfig.FakeResource{
				Annotations: config.Annotations{
					Annotations: config.AnnotationFields{Description: "foo"},
				},
			},
			expected: false,
		},
		{
			opts: listOptions{tags: []string{"one"}},
			resource: &testconfig.FakeResource{
				Annotations: config.Annotations{
					Annotations: config.AnnotationFields{
						Tags: []string{"one", "two"},
					},
				},
			},
			expected: true,
		},
	}

	for _, testcase := range testcases {
		actual := include(testcase.resource, testcase.opts)
		assert.Check(t, is.Equal(testcase.expected, actual))
	}
}
