package generator

import (
	"testing"

	"github.com/labd/mach-composer/internal/config"
	"github.com/labd/mach-composer/internal/model"
	"github.com/labd/mach-composer/internal/model/commercetools"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	cfg := config.Config{
		MACHComposer: config.MACHComposer{
			Version: "1.0.0",
		},
		Global: config.Global{
			Environment: "test",
			Cloud:       "aws",
			TerraformConfig: model.TerraformConfig{
				AwsRemoteState: &model.AWSTFState{
					Bucket:    "your bucket",
					KeyPrefix: "mach",
					Region:    "eu-central-1",
				},
			},
		},
		Sites: []model.Site{
			{
				Name:       "",
				Identifier: "my-site",
				RawEndpoints: map[string]any{
					"main": "api.my-site.nl",
					"internal": map[string]any{
						"throttling_burst_limit": 5000,
						"throttling_rate_limit":  10000,
						"url":                    "internal-api.my-site.nl",
					},
				},
				Commercetools: &commercetools.Settings{
					ProjectKey:   "my-site",
					ClientID:     "<client-id>",
					ClientSecret: "<client-secret>",
					Scopes:       "manage_api_clients:my-site manage_project:my-site view_api_clients:my-site",
					ProjectSettings: &commercetools.ProjectSettings{
						Languages:  []string{"en-GB", "nl-NL"},
						Currencies: []string{"GBP", "EUR"},
						Countries:  []string{"GB", "NL"},
					},
				},
				Components: []model.SiteComponent{
					{
						Name: "your-component",
						Variables: map[string]any{
							"FOO_VAR": "my-value",
						},
						Secrets: map[string]any{
							"MY_SECRET": "secretvalue",
						},
					},
				},
				AWS: &model.SiteAWS{
					AccountID: "123456789",
					Region:    "eu-central-1",
				},
			},
		},
		Components: []model.Component{
			{
				Name:         "your-component",
				Source:       "git::https://github.com/<username>/<your-component>.git//terraform",
				Version:      "0.1.0",
				Integrations: []string{"aws", "commercetools"},
			},
		},
	}
	config.Process(&cfg)
	body, err := Render(&cfg, &cfg.Sites[0])
	assert.NoError(t, err)
	assert.NotEmpty(t, body)
}
