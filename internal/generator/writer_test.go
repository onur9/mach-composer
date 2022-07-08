package generator

import (
	"testing"

	"github.com/labd/mach-composer/internal/config"
	"github.com/labd/mach-composer/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestFileLocations(t *testing.T) {
	cfg := &config.Config{
		Sites: []model.Site{
			{
				Identifier: "my-site",
			},
		},
	}

	options := &GenerateOptions{
		OutputPath: "deployments/something",
	}

	actual := FileLocations(cfg, options)
	expected := map[string]string{
		"my-site": "deployments/something/my-site",
	}

	assert.EqualValues(t, expected, actual)
}
