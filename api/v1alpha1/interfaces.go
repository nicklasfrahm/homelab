package v1alpha1

// Validatable is an interface for validating a resource.
type Validatable interface {
	// Validate validates the resource.
	Validate() error
}
