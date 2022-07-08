package main

import (
	"context"

	"github.com/labd/mach-composer/internal/generator"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the Terraform files.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := generateFunc(context.Background(), args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	registerGenerateFlags(generateCmd)
}

func generateFunc(ctx context.Context, args []string) error {
	genOptions := &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	}

	configs := LoadConfigs(ctx)
	for _, filename := range generateFlags.fileNames {
		cfg := configs[filename]

		_, err := generator.WriteFiles(cfg, genOptions)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
