package errors

import "errors"

var (
	// As finds the first error in err's chain that matches target, and if so, sets target to that error value and returns true. Otherwise, it returns false.
	As = errors.As

	// Is reports whether any error in err's chain matches target.
	Is = errors.Is

	// New returns an error that formats as the given text. Each call to New returns a distinct error value even if the text is identical.
	New = errors.New

	// Unwrap returns the result of calling the Unwrap method on err, if err's type contains an Unwrap method returning error. Otherwise, Unwrap returns nil.
	Unwrap = errors.Unwrap
)
