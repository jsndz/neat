package utils

import (
	"fmt"
	"time"
)

func AuthorLine() string {
	now := time.Now()
	_, offset := now.Zone()

	// convert offset seconds to +0530 format
	tz := fmt.Sprintf("%+03d%02d", offset/3600, (offset%3600)/60)

	return fmt.Sprintf(
		"author neat <neat@mail.com> %d %s",
		now.Unix(),
		tz,
	)
}

func CommitterLine() string {
	now := time.Now()
	_, offset := now.Zone()
	tz := fmt.Sprintf("%+03d%02d", offset/3600, (offset%3600)/60)

	return fmt.Sprintf(
		"committer neat <neat@mail.com> %d %s",
		now.Unix(),
		tz,
	)
}
