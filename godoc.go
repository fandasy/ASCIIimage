// Package ASCIIImage serves as the root package providing version information and documentation.
//
// Note: All implementation logic resides in subpackages:
//   - ./core/ : Low-level ASCII generation and rendering
//   - ./api/  : High-level client for common use cases
//
// This package exists primarily to:
//   - Provide unified documentation
//   - Track version information
//   - Serve as the main import path
package ASCIIImage

// Info contains package metadata for version tracking and documentation.
// This struct is intentionally simple as the package contains no implementation.
type Info struct {
	Name    string
	Version string
}

// PackageInfo returns metadata about the package distribution.
//
// This is a documentation-focused function that:
//   - Provides version information for dependency checking
//   - Serves as a root documentation anchor
//   - Contains no business logic
//
// For actual functionality, import:
//   - ./api for the high-level client
//   - ./core for low-level operations
func PackageInfo() Info {
	return Info{
		Name:    "github.com/fandasy/ASCIIimage/v2",
		Version: "v2.2.3",
	}
}
