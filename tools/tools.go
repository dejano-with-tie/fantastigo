//go:build tools
// +build tools

// tools package implements _Tools as Dependencies_ paradigm.
// This package could be converted to a module
// Good reads:
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://marcofranssen.nl/manage-go-tools-via-go-modules
package tools

import (
	_ "github.com/deepmap/oapi-codegen/pkg/codegen" // OpenAPI 3 Code Generator
)
