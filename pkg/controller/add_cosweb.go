package controller

import (
	"cosweb-operator/pkg/controller/cosweb"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, cosweb.Add)
}
