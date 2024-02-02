package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/desmos-labs/athena/v2/x/profiles-score/scorers/twitter"
	"github.com/desmos-labs/athena/v2/x/profiles-score/scorers/utils"
)

func TestUnmarshalConfig(t *testing.T) {
	testCases := []struct {
		name      string
		config    string
		nodeName  string
		value     interface{}
		shouldErr bool
		expFound  bool
		check     func(value interface{})
	}{
		{
			name:      "missing top level value returns properly",
			config:    ``,
			value:     map[string]string{},
			shouldErr: false,
			expFound:  false,
		},
		{
			name: "not found config returns properly",
			config: `
scorers:
`,
			nodeName:  "github",
			value:     map[string]string{},
			shouldErr: false,
			expFound:  false,
		},
		{
			name: "found config returns properly",
			config: `
scorers:
  twitter:
    token: "custom_token"
`,
			nodeName:  "twitter",
			value:     &twitter.Config{},
			shouldErr: false,
			expFound:  true,
			check: func(value interface{}) {
				cfg, ok := value.(*twitter.Config)
				require.True(t, ok)
				require.Equal(t, cfg.Token, "custom_token")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			found, err := utils.UnmarshalConfig([]byte(tc.config), tc.nodeName, tc.value)
			if tc.shouldErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expFound, found)
				if tc.expFound {
					tc.check(tc.value)
				}
			}
		})
	}
}
