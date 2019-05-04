package connect

import (
	"path/filepath"
	"strings"
)

// ActivityFormat is a file format for importing and exporting activities.
type ActivityFormat int

const (
	// ActivityFormatFIT is the "original" Garmin format.
	ActivityFormatFIT ActivityFormat = iota

	// ActivityFormatTCX is Training Center XML (TCX) format.
	ActivityFormatTCX

	// ActivityFormatGPX will export as GPX - the GPS Exchange Format.
	ActivityFormatGPX

	// ActivityFormatKML will export KML files compatible with Google Earth.
	ActivityFormatKML

	// ActivityFormatCSV will export splits as CSV.
	ActivityFormatCSV

	activityFormatMax
	activityFormatInvalid
)

const (
	// ErrUnknownFormat will be returned if the activity file format is unknown.
	ErrUnknownFormat = Error("Unknown format")
)

var (
	activityFormatTable = map[string]ActivityFormat{
		"fit": ActivityFormatFIT,
		"tcx": ActivityFormatTCX,
		"gpx": ActivityFormatGPX,
		"kml": ActivityFormatKML,
		"csv": ActivityFormatCSV,
	}
)

// Extension returns an appropriate filename extension for format.
func (f ActivityFormat) Extension() string {
	for extension, format := range activityFormatTable {
		if format == f {
			return extension
		}
	}

	return ""
}

// FormatFromExtension tries to guess the format from a file extension.
func FormatFromExtension(extension string) (ActivityFormat, error) {
	extension = strings.ToLower(extension)

	format, found := activityFormatTable[extension]
	if !found {
		return activityFormatInvalid, ErrUnknownFormat
	}

	return format, nil
}

// FormatFromFilename tries to guess the format based on a filename (or path).
func FormatFromFilename(filename string) (ActivityFormat, error) {
	extension := filepath.Ext(filename)
	extension = strings.TrimPrefix(extension, ".")

	return FormatFromExtension(extension)
}
