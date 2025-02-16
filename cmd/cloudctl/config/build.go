package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	cloud "github.com/nicklasfrahm/cloud/api/v1beta1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
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

			schemas := map[string]bool{
				"machines": true,
			}

			srcDir := args[0]
			dstDir := args[1]

			fakeServer := NewFakeServer()
			fakeServer.Register(cloud.SchemeBuilder)

			// Generate documents.
			if err := fakeServer.Generate(dstDir); err != nil {
				return fmt.Errorf("failed to generate documents: %w", err)
			}

			apiGroupDir := path.Join(dstDir, "apis", cloud.SchemeBuilder.GroupVersion.Group, cloud.SchemeBuilder.GroupVersion.Version)

			for schema := range schemas {
				srcSchemaDir := path.Join(srcDir, schema)
				dstSchemaDir := path.Join(apiGroupDir, schema)

				if err := os.MkdirAll(dstSchemaDir, 0755); err != nil {
					return fmt.Errorf("failed to create schema directory: %w", err)
				}

				// Read files.
				entries, err := os.ReadDir(path.Join(srcDir, schema))
				if err != nil {
					return fmt.Errorf("failed to read schema directory: %w", err)
				}

				for _, entry := range entries {
					// We do not expect subdirectories.
					if entry.IsDir() {
						continue
					}

					var err error
					switch schema {
					case "machines":
						err = transcodeManifest[cloud.Machine](srcSchemaDir, dstSchemaDir, entry.Name())
					}

					if err != nil {
						return fmt.Errorf("failed to transcode manifest: %w", err)
					}
				}
			}

			return nil
		},
	}

	return cmd
}

// FakeServer is a fake Kubernetes API server.
type FakeServer struct {
	RootPaths 	metav1.RootPaths
	APIGroups 	metav1.APIGroupList
}

// NewFakeServer returns a new fake server.
func NewFakeServer() *FakeServer {
	return &FakeServer{
		RootPaths: metav1.RootPaths{
			Paths: []string{
				"/apis",
				"/apis/",
			},
		},
	}
}

// RegisterGroupVersion registers an API via a scheme builder.
func (s *FakeServer) Register(schemeBuilder *scheme.Builder) {
	apiGroupPath := fmt.Sprintf("/apis/%s", schemeBuilder.GroupVersion.Group)
	apiGroupVersionPath := fmt.Sprintf("/apis/%s/%s", schemeBuilder.GroupVersion.Group, schemeBuilder.GroupVersion.Version)

	s.RootPaths.Paths = append(s.RootPaths.Paths, apiGroupPath, apiGroupVersionPath)

	s.APIGroups.Groups = append(s.APIGroups.Groups, metav1.APIGroup{
		Name: schemeBuilder.GroupVersion.Group,
		Versions: []metav1.GroupVersionForDiscovery{
			{
				GroupVersion: schemeBuilder.GroupVersion.String(),
				Version:      schemeBuilder.GroupVersion.Version,
			},
		},
		PreferredVersion: metav1.GroupVersionForDiscovery{
			GroupVersion: schemeBuilder.GroupVersion.String(),
			Version:      schemeBuilder.GroupVersion.Version,
		},
	})
}

// generateIndex generates the index file.
func (s *FakeServer) generateIndex(dir string, data interface{}) error {
	// Create target directory.
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	indexFile := path.Join(dir, "index.json")
	if err := os.WriteFile(indexFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write index file: %w", err)
	}

	return nil
}

// Generate generates the documents to be served by static file server.
func (s *FakeServer) Generate(dstDir string) error {
	// Remove any existing target directory.
	if err := os.RemoveAll(dstDir); err != nil {
		// It's okay if the target directory does not exist yet,
		// we just need to make sure it's empty.
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove target directory: %w", err)
		}
	}

	// Recreate target directory.
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Generate root paths.
	if err := s.generateIndex(dstDir, s.RootPaths); err != nil {
		return fmt.Errorf("failed to generate root paths: %w", err)
	}

	// Generate API groups.
	if err := s.generateIndex(path.Join(dstDir, "apis"), s.APIGroups); err != nil {
		return fmt.Errorf("failed to generate API groups: %w", err)
	}

	return nil
}

// ManifestDecoder decodes a Kubernetes manifest.
type ManifestDecoder[T any] struct {
	reader  io.Reader
	decoder runtime.Decoder
}

// NewManifestDecoder creates a new manifest decoder.
func NewManifestDecoder[T any](reader io.Reader) *ManifestDecoder[T] {
	return &ManifestDecoder[T]{
		reader: reader,
		decoder: yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme),
	}
}

// Decode decodes a Kubernetes manifest.
func (d *ManifestDecoder[T]) Decode(entity *T) (error) {
	// Create a buffer from the reader.
	buf := bytes.NewBuffer(nil)

	if _, err := io.Copy(buf, d.reader); err != nil {
		return fmt.Errorf("failed to copy reader to buffer: %w", err)
	}

	obj, _, err := d.decoder.Decode(buf.Bytes(), nil, nil)
	if err != nil {
		return fmt.Errorf("failed to decode Kubernetes manifest: %w", err)
	}

	err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.(*unstructured.Unstructured).UnstructuredContent(), entity)
	if err != nil {
		return fmt.Errorf("failed to convert Kubernetes manifest: %w", err)
	}

	return nil
}

// transcodeManifest transcodes a Kubernetes manifest.
func transcodeManifest[T any](srcDir, dstDir, name string) error {
	// Open file.
	srcPath := path.Join(srcDir, name)
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open schema file: %w", err)
	}
	defer src.Close()

	// Decode file.
	decoder := NewManifestDecoder[T](src)

	var obj T
	if err := decoder.Decode(&obj); err != nil {
		return fmt.Errorf("failed to decode schema file: %w", err)
	}

	// Encode file.
	dstPath := path.Join(dstDir, strings.TrimSuffix(name, path.Ext(name)) + ".json")
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create schema file: %w", err)
	}

	encoder := json.NewEncoder(dstFile)
	if err := encoder.Encode(obj); err != nil {
		return fmt.Errorf("failed to encode schema file: %w", err)
	}

	// Close file.
	if err := dstFile.Close(); err != nil {
		return fmt.Errorf("failed to close schema file: %w", err)
	}

	return nil
}
