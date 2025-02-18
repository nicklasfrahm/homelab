package config

import (
	"fmt"
	"os"
	"path"

	cloud "github.com/nicklasfrahm/cloud/api/v1beta1"
	"github.com/nicklasfrahm/cloud/pkg/kubeenc"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

// BuildCommand returns the build command.
func BuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build <src_dir> <dst_dir>",
		Short: "Build configuration into static files",
		Long: `Build configuration into static files
that can be served by a web server.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("expected exactly two arguments")
			}

			inputDir := args[0]
			outputDir := args[1]

			repository := NewConfigRepository()

			schemas := map[string]ResourceLoader{
				"machines": Load(&repository.Machines.Items),
			}

			for schema, load := range schemas {
				schemaDir := path.Join(inputDir, schema)

				if err := load(schemaDir); err != nil {
					return fmt.Errorf("failed to load schema: %w", err)
				}
			}

			versionDir := path.Join(outputDir, cloud.GroupVersion.Version)
			if err := repository.Build(versionDir); err != nil {
				return fmt.Errorf("failed to build configuration: %w", err)
			}

			return nil
		},
	}

	return cmd
}


// ConfigRepository is a configuration repository.
type ConfigRepository struct {
	Machines cloud.MachineList
}

// NewConfigRepository creates a new configuration repository.
func NewConfigRepository() *ConfigRepository {
	return &ConfigRepository{
		Machines: cloud.MachineList{
			Items: []cloud.Machine{},
		},
	}
}

// ResourceLoader is a function that loads a resource into a repository.
type ResourceLoader func(srcDir string) error

// Load loads the configuration. This is not optimized for performance,
// we should most likely make this concurrent.
func Load[T any](repository *[]T) ResourceLoader {
	return func(schemaDir string) error {
		// Read files.
		entries, err := os.ReadDir(schemaDir)
		if err != nil {
			return fmt.Errorf("failed to read schema directory: %w", err)
		}

		for _, entry := range entries {
			// We do not expect subdirectories.
			if entry.IsDir() {
				continue
			}

			resourceManifest := path.Join(schemaDir, entry.Name())

			rawResource, err := os.ReadFile(resourceManifest)
			if err != nil {
				return fmt.Errorf("failed to read resource: %w", err)
			}

			decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

			obj, _, err := decoder.Decode(rawResource, nil, nil)
			if err != nil {
				return fmt.Errorf("failed to decode resource manifest: %w", err)
			}

			entity := new(T)

			err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).UnstructuredContent(), entity)
			if err != nil {
				return fmt.Errorf("failed to convert Kubernetes manifest: %w", err)
			}

			*repository = append(*repository, *entity)
		}

		return nil
	}
}

// Build builds a resource into a file.
func Build[T runtime.Object](dstFile string, resource T) error {
	if err := os.MkdirAll(path.Dir(dstFile), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	file, err := os.Create(dstFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := kubeenc.NewJSONEncoder(file)

	cloudScheme, err := cloud.SchemeBuilder.Build()
	if err != nil {
		return fmt.Errorf("failed to build scheme: %w", err)
	}

	data, err := encoder.EncodeWithScheme(resource, cloudScheme)
	if err != nil {
		return fmt.Errorf("failed to encode resource: %w", err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Build builds the configuration repository into static files.
func (r *ConfigRepository) Build(dstDir string) error {
	// Clear the destination directory.
	if err := os.RemoveAll(dstDir); err != nil {
		// Ignore errors if the directory does not exist.
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove destination directory: %w", err)
		}
	}

	// Create the destination directory.
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	machineIndex := path.Join(dstDir, "machines", "index.json")
	if err := Build(machineIndex, &r.Machines); err != nil {
		return fmt.Errorf("failed to build machine index: %w", err)
	}

	for _, machine := range r.Machines.Items {
		machineFile := path.Join(dstDir, "machines", machine.Name + ".json")
		if err := Build(machineFile, &machine); err != nil {
			return fmt.Errorf("failed to build machine: %w", err)
		}
	}

	return nil
}
