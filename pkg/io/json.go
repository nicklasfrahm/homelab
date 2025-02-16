package io

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
	printer, err := genericclioptions.
		NewPrintFlags("render").
		WithTypeSetter(scheme.Scheme).
		WithDefaultOutput("json").
		ToPrinter()
	if err != nil {
		return nil, fmt.Errorf("failed to create printer: %w", err)
	}

	if err := printer.PrintObj(obj, e.writer); err != nil {
		return nil, fmt.Errorf("failed to print secret: %w", err)
	}

	return nil, nil
}
