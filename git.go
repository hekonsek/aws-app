package awsom

import (
	"fmt"
	"strconv"
	"strings"
)

func CurrentVersion(workingDir string) (string, error) {
	tagsOut, err := Exec{
		Command:    "git tag -l",
		WorkingDir: workingDir,
	}.Run()
	if err != nil {
		return "", err
	}

	var latest int64 = 0
	for _, tag := range tagsOut {
		versionNumber, err := strconv.ParseInt(strings.Split(tag, ".")[1], 0, 64)
		if err != nil {
			continue
		}
		if versionNumber > latest {
			latest = versionNumber
		}
	}

	return fmt.Sprintf("0.%d", latest), nil
}

func NextVersion(workingDir string) (string, error) {
	currentVersion, err := CurrentVersion(workingDir)
	if err != nil {
		return "", err
	}
	versionNumber, err := strconv.ParseInt(strings.Split(currentVersion, ".")[1], 0, 64)
	if err != nil {
		return "", err

	}
	return fmt.Sprintf("0.%d", versionNumber+1), nil
}
