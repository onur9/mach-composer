package main

import (
	"context"

	"github.com/labd/mach-composer/internal/generator"
	"github.com/labd/mach-composer/internal/runner"
	"github.com/spf13/cobra"
)

var planFlags struct {
	reuse      bool
	components []string
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plan the configuration.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if err := planFunc(context.Background(), args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	registerGenerateFlags(planCmd)
	planCmd.Flags().BoolVarP(&planFlags.reuse, "reuse", "", false, "Suppress a terraform init for improved speed (not recommended for production usage)")
	planCmd.Flags().StringArrayVarP(&planFlags.components, "component", "c", []string{}, "")
}

func planFunc(ctx context.Context, args []string) error {
	configs := LoadConfigs(ctx)
	allPaths := make(map[string]map[string]string)

	// Write the generate files for each config
	genOptions := &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	}
	for _, filename := range generateFlags.fileNames {
		cfg := configs[filename]
		paths, err := generator.WriteFiles(cfg, genOptions)
		if err != nil {
			panic(err)
		}
		allPaths[filename] = paths
	}

	// Plan the generate files
	options := &runner.PlanOptions{
		Reuse: planFlags.reuse,
		Site:  generateFlags.siteName,
	}
	for _, filename := range generateFlags.fileNames {
		cfg := configs[filename]
		paths := allPaths[filename]
		runner.TerraformPlan(cfg, paths, options)
	}

	return nil
}
