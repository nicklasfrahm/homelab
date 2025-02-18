package kubeenc

import (
	"fmt"
	"io"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/scheme"
)

// JSONEncoder is a custom JSON encoder.
type JSONEncoder struct {
	writer io.Writer
}

// NewJSONEncoder creates a new JSON encoder.
func NewJSONEncoder(writer io.Writer) *JSONEncoder {
	return &JSONEncoder{
		writer: writer,
	}
}

// Encode encodes JSON for a Kubernetes API object.
func (e *JSONEncoder) Encode(obj runtime.Object) ([]byte, error) {
	return e.EncodeWithScheme(obj, scheme.Scheme)
}

// EncodeWithScheme encodes JSON for a Kubernetes API object with a custom scheme.
func (e *JSONEncoder) EncodeWithScheme(obj runtime.Object, customScheme *runtime.Scheme) ([]byte, error) {
	printer, err := genericclioptions.
		NewPrintFlags("render").
		WithTypeSetter(customScheme).
		WithDefaultOutput("json").
		ToPrinter()
	if err != nil {
		return nil, fmt.Errorf("failed to create printer: %w", err)
	}

	if err := printer.PrintObj(obj, e.writer); err != nil {
		return nil, fmt.Errorf("failed to print object: %w", err)
	}

	return nil, nil
}
