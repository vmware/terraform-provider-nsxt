package util

import (
	"log"

	"github.com/hashicorp/go-version"
)

var NsxVersion = ""

func NsxVersionLower(ver string) bool {

	requestedVersion, err1 := version.NewVersion(ver)
	currentVersion, err2 := version.NewVersion(NsxVersion)
	if err1 != nil || err2 != nil {
		log.Printf("[ERROR] Failed perform version check for version %s", ver)
		return true
	}
	return currentVersion.LessThan(requestedVersion)
}

func NsxVersionHigherOrEqual(ver string) bool {

	requestedVersion, err1 := version.NewVersion(ver)
	currentVersion, err2 := version.NewVersion(NsxVersion)
	if err1 != nil || err2 != nil {
		log.Printf("[ERROR] Failed perform version check for version %s", ver)
		return false
	}
	return currentVersion.Compare(requestedVersion) >= 0
}
