package observable

// Annotation represents a key-value pair used for annotating observer values.
type Annotation struct {
	Key   string
	Value string
}

// Annotations is a map that associates string keys with Annotation values.
// It allows the storage of multiple annotations for observer values.
type Annotations map[string]Annotation
