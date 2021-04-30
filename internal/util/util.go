// Package util provides local helper functions that
// do not interact with any external services
package util

import (
	"github.com/walkamongus/card-search/internal/hsapi"
)

// Contains checks whether a slice contains a specific number
func Contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// GetRarityName searches a list of Rarity object metadata
// and looks up the rarity name by a provided ID
func GetRarityName(id int64, metadata []hsapi.Rarity) string {
	for _, m := range metadata {
		if m.ID == id {
			return m.Name
		}
	}
	return ""
}

// GetClassName searches a list of Class object metadata
// and looks up the class name by a provided ID
func GetClassName(id int64, metadata []hsapi.Class) string {
	for _, m := range metadata {
		if m.ID == id {
			return m.Name
		}
	}
	return ""
}

// GetName searches a list of GenericMetadata object metadata
// and looks up the object name by a provided ID
func GetName(id int64, metadata []hsapi.GenericMetadata) string {
	for _, m := range metadata {
		if m.ID == id {
			return m.Name
		}
	}
	return ""
}

// GetSetName searches a list of Set object metadata
// and looks up the object name by a provided ID.
// Both alias set IDs and top-level IDs are searched.
func GetSetName(id int64, metadata []hsapi.Set) string {
	for _, m := range metadata {
		if Contains(m.AliasSetIds, id) {
			return m.Name
		}
		if m.ID == id {
			return m.Name
		}
	}
	return ""
}
