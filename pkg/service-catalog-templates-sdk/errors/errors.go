// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package errors

import "errors"

const (
	// ErrorUnmanagedResource means the resource exists, but is not owned by
	// its corresponding shadow resource.
	ErrorUnmanagedResource = "resource exists, but is not managed by a templated resource"
)

func NewUnmanagedResource() error {
	return errors.New(ErrorUnmanagedResource)
}

// IsUnmanagedResource returns true if the specified error was created by NewUnmanagedResource.
func IsUnmanagedResource(err error) bool {
	return err.Error() == ErrorUnmanagedResource
}
