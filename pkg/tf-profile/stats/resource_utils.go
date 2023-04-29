package tfprofile

import (
	"strings"
)

// Extract the top level module.
// For example, "module.mymod.aws_subnet.test" will return "module.mymod"
func getTopLevelModule(name string) string {
	split := strings.Split(name, ".")
	if len(split) < 2 {
		return ""
	}
	if split[0] == "module" {
		return split[0] + "." + split[1]
	}
	return ""
}

// Return the amount of nested modules for a resource.
// E.g. aws_subnet.test => 0
// E.g. module.mod1.aws_subnet.test => 1.
// E.g. module.mod1.module.mod2.aws_subnet.test => 2.
func getModuleDepth(name string) int {
	tokens := len(strings.Split(name, "."))
	return (tokens - 2) / 2
}

// Given a full resource name, return the name of the deepest module it belongs to
// (including parent modules)
func getModule(name string) string {
	split := strings.Split(name, ".")
	return strings.Join(split[:len(split)-2], ".")
}

// Given a full resource name, return only the name of the deepest module
// without parent modules
func getLeafModuleName(name string) string {
	split := strings.Split(name, ".")
	leaf := ""

	for idx, s := range split {
		if s == "module" {
			leaf = split[idx+1]
		}
	}
	return leaf
}
