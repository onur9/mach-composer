package generator

import (
	"embed"
	"fmt"
	"strings"

	"github.com/flosch/pongo2/v5"
	"github.com/labd/mach-composer/config"
)

//go:embed templates/*
var templates embed.FS

type TemplateRenderer struct {
	templateSet *pongo2.TemplateSet

	servicesTemplate  *pongo2.Template
	componentTemplate *pongo2.Template
	tfConfigTemplate  *pongo2.Template
}

var renderer TemplateRenderer

func init() {
	registerFilters()

	renderer.templateSet = pongo2.NewSet("", &EmbedLoader{Content: templates})
	renderer.servicesTemplate = pongo2.Must(renderer.templateSet.FromFile("services.tf"))
	renderer.componentTemplate = pongo2.Must(renderer.templateSet.FromFile("component.tf"))
	renderer.tfConfigTemplate = pongo2.Must(renderer.templateSet.FromFile("config.tf"))
}

func Render(cfg *config.MachConfig, site *config.Site) (string, error) {
	result := []string{
		"# This file is auto-generated by MACH composer",
		fmt.Sprintf("# Site: %s", site.Identifier),
	}

	if val, err := RenderTerraformConfig(cfg, site); err == nil {
		result = append(result, val)
	} else {
		return "", err
	}

	if val, err := RenderServices(cfg, site); err == nil {
		result = append(result, val)
	} else {
		return "", err
	}

	// Add components
	for i := range site.Components {
		if val, err := RenderComponent(cfg, site, &site.Components[i]); err == nil {
			result = append(result, val)
		} else {
			return "", err
		}
	}

	content := strings.Join(result, "\n")
	return content, nil
}

func RenderTerraformConfig(cfg *config.MachConfig, site *config.Site) (string, error) {
	return renderer.tfConfigTemplate.Execute(pongo2.Context{
		"global":    cfg.Global,
		"site":      site,
		"variables": cfg.Variables,
	})
}

func RenderServices(cfg *config.MachConfig, site *config.Site) (string, error) {
	return renderer.servicesTemplate.Execute(pongo2.Context{
		"global":    cfg.Global,
		"site":      site,
		"variables": cfg.Variables,
	})
}

func RenderComponent(cfg *config.MachConfig, site *config.Site, component *config.SiteComponent) (string, error) {
	return renderer.componentTemplate.Execute(pongo2.Context{
		"global":     cfg.Global,
		"site":       site,
		"variables":  cfg.Variables,
		"component":  component,
		"definition": component.Definition,
	})
}
