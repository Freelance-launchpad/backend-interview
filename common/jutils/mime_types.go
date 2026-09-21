package jutils

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/gabriel-vasile/mimetype"
)

func isCSV(data []byte) bool {
	potentialSeparators := []rune{';', '|'}

	for _, sep := range potentialSeparators {
		r := csv.NewReader(bytes.NewReader(data))
		r.Comma = sep
		r.Comment = '#'

		headers, err := r.Read()
		if err != nil {
			continue
		}
		numCols := len(headers)

		if numCols < 2 {
			continue
		}

		records, err := r.Read()
		if errors.Is(err, io.EOF) { // ensure at least 2 lines
			continue
		}
		if err != nil || len(records) != numCols {
			continue
		}

		return true
	}

	return false
}

func GetMimeType(rs io.ReadSeeker) (_ string, err error) {
	defer jerror.Wrap(&err)
	defer func() {
		if cursor, seekErr := rs.Seek(0, io.SeekStart); seekErr != nil {
			err = errors.Join(err, seekErr)
		} else if cursor != 0 {
			err = errors.Join(err, fmt.Errorf("unable to seek 0"))
		}
	}()

	if cursor, err := rs.Seek(0, io.SeekStart); err != nil {
		return "", err
	} else if cursor != 0 {
		return "", fmt.Errorf("unable to seek 0")
	}

	buffer := make([]byte, 3072)
	n, err := rs.Read(buffer)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("failed to read from reader: %w", err)
	}
	buffer = buffer[:n]

	mime := mimetype.Detect(buffer)
	if mime.Is("text/plain") {
		if isCSV(buffer) {
			mime = mimetype.Lookup("text/csv")
		}
	}

	return mime.String(), nil
}
