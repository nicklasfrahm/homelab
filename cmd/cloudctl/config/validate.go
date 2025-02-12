package config

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	cloud "github.com/nicklasfrahm/cloud/api/v1beta1"
)

// ValidateCommand returns the validate command.
func ValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <directory>",
		Short: "Validate manifests in a directory",
		Long:  `Validate manifests in a directory.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("expected exactly one argument")
			}

			schemas := map[string]bool{
				"machines": true,
			}

			// Read folders in directory and check if they are in the schemas map.
			configDir := args[0]
			entries, err := os.ReadDir(configDir)
			if err != nil {
				return fmt.Errorf("failed to read directory: %w", err)
			}

			for _, entry := range entries {
				// Skip files that are directly in the config directory.
				if !entry.IsDir() {
					continue
				}

				// Check if the folder is in the schemas map.
				if _, ok := schemas[entry.Name()]; !ok {
					fmt.Printf("🟡 Skipping %s: not in schemas\n", entry.Name())
					continue
				}

				// Validate the schema of the directory.
				fmt.Printf("🟢 Validating schema: %s\n", entry.Name())

				schemaDir := path.Join(configDir, entry.Name())

				switch entry.Name() {
				case "machines":
					if err := validateSchema[cloud.Machine](schemaDir); err != nil {
						return fmt.Errorf("failed to validate schema: %w", err)
					}
				}
			}

			return nil
		},
	}
}

// validateSchema validates the schema of a directory.
func validateSchema[T any](directory string) error {
	// Read files in directory and check if they match the schema.
	entries, err := os.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		// Skip directories.
		if entry.IsDir() {
			continue
		}

		// Read file content and unmarshal it into the schema.
		filePath := path.Join(directory, entry.Name())
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		instance := new(T)

		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&instance); err != nil {
			fmt.Printf("🔴 >> %s: %v\n", entry.Name(), errors.Unwrap(err))

			continue
		}

		fmt.Printf("🟢 >> %s\n", entry.Name())
	}

	return nil
}
