package flag

import (
	"flag"
	"testing"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestCli_Read(t *testing.T) {
	tests := map[string]struct {
		givenExistingConfig map[string]interface{}
		givenFlags          []cli.Flag
		givenArgs           []string
		givenUseDefault     bool
		givenAliases        map[string]string
		expectedErr         string
		expectedConfig      map[string]interface{}
	}{
		"GivenNoExistingConfig_WhenStringFlagInRoot_ThenParseFlag": {
			givenExistingConfig: map[string]interface{}{},
			givenFlags: []cli.Flag{
				&cli.StringFlag{Name: "flag"},
			},
			givenArgs: []string{"--flag", "foo"},
			expectedConfig: map[string]interface{}{
				"flag": "foo",
			},
		},
		"GivenNoExistingConfig_WhenStringFlagNested_ThenParseFlag": {
			givenExistingConfig: map[string]interface{}{},
			givenFlags: []cli.Flag{
				&cli.StringFlag{Name: "nested-flag"},
			},
			givenArgs: []string{"--nested-flag", "foo"},
			expectedConfig: map[string]interface{}{
				"nested": map[string]interface{}{
					"flag": "foo",
				},
			},
		},
		"GivenExistingConfig_WhenStringFlagInRoot_ThenOverwriteExisting": {
			givenExistingConfig: map[string]interface{}{
				"flag": "bar",
			},
			givenFlags: []cli.Flag{
				&cli.StringFlag{Name: "flag"},
			},
			givenArgs: []string{"--flag", "foo"},
			expectedConfig: map[string]interface{}{
				"flag": "foo",
			},
		},
		"GivenExistingConfig_WhenStringFlagNested_ThenOverwriteExisting": {
			givenExistingConfig: map[string]interface{}{
				"nested": map[string]interface{}{
					"flag": "bar",
				},
			},
			givenFlags: []cli.Flag{
				&cli.StringFlag{Name: "nested-flag"},
			},
			givenArgs: []string{"--nested-flag", "foo"},
			expectedConfig: map[string]interface{}{
				"nested": map[string]interface{}{
					"flag": "foo",
				},
			},
		},
		"GivenExistingConfig_WhenNoArgsGiven_ThenIgnoreDefaultValue": {
			givenExistingConfig: map[string]interface{}{
				"nested": map[string]interface{}{
					"flag": "bar",
				},
			},
			givenFlags: []cli.Flag{
				&cli.StringFlag{Name: "nested-flag", Value: "foo", Aliases: []string{"alias"}},
			},
			givenArgs:      []string{},
			expectedConfig: map[string]interface{}{},
		},
		"GivenAliases_WhenStringFlagNested_ThenCreateNestedConfig": {
			givenFlags: []cli.Flag{
				&cli.StringFlag{Name: "flag"},
			},
			givenAliases: map[string]string{
				"flag": "nested-flag",
			},
			givenArgs: []string{"--flag", "foo"},
			expectedConfig: map[string]interface{}{
				"nested": map[string]interface{}{
					"flag": "foo",
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fs := flag.NewFlagSet("", flag.PanicOnError)
			for _, fl := range tt.givenFlags {
				strFlag := fl.(*cli.StringFlag)
				fs.String(strFlag.Name, strFlag.Value, strFlag.Usage)
			}

			require.NoError(t, fs.Parse(tt.givenArgs))
			ctx := cli.NewContext(&cli.App{}, fs, nil)

			k := koanf.New(".")
			require.NoError(t, k.Load(confmap.Provider(tt.givenExistingConfig, k.Delim()), nil))

			p := Provider(ctx, "-", k, tt.givenAliases)

			result, err := p.Read()
			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedConfig, result)
		})
	}
}
