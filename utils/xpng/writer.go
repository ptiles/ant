package xpng

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"image"
	"image/png"
	"io"
)

func Encode(w io.Writer, m image.Image, txt map[string]string) error {
	if len(txt) == 0 {
		return png.Encode(w, m)
	}

	return png.Encode(&txtWriter{w: w, txt: txt}, m)
}

type txtWriter struct {
	w   io.Writer
	txt map[string]string

	idatFound bool
}

func (tw *txtWriter) Write(b []byte) (int, error) {
	if !tw.idatFound && len(b) >= 8 && b[4] == 'I' && b[5] == 'D' && b[6] == 'A' && b[7] == 'T' {
		tw.idatFound = true

		for k, v := range tw.txt {
			n, err := tw.writeItxtChunk(k, v)
			if err != nil {
				return n, err
			}
		}

		n, err := tw.writeItxtChunk("Software", "github.com/ptiles/ant")
		if err != nil {
			return n, err
		}
	}

	return tw.w.Write(b)
}

func (tw *txtWriter) writeItxtChunk(key, value string) (int, error) {
	itxtLength := len(key) + 5 + len(value)
	if itxtLength > maxTxtSize {
		return 0, errors.New("iTXt chunk is too large")
	}

	// 4 bytes total length; 4 bytes chunk name; 4 bytes checksum
	chunk := make([]byte, 0, 4+4+itxtLength+4)

	// Length
	chunk = binary.BigEndian.AppendUint32(chunk, uint32(itxtLength))

	// Chunk name
	chunk = append(chunk, 'i', 'T', 'X', 't')

	// Keyword:             1-79 bytes (character string)
	chunk = append(chunk, key...)
	// Null separator:      1 byte
	chunk = append(chunk, 0)

	// Compression flag:    1 byte
	chunk = append(chunk, 0)

	// Compression method:  1 byte
	chunk = append(chunk, 0)

	// Language tag:        0 or more bytes (character string)
	// Null separator:      1 byte
	chunk = append(chunk, 0)

	// Translated keyword:  0 or more bytes
	// Null separator:      1 byte
	chunk = append(chunk, 0)

	// Text:                0 or more bytes
	chunk = append(chunk, value...)

	crc := crc32.ChecksumIEEE(chunk[4:])
	chunk = binary.BigEndian.AppendUint32(chunk, crc)

	return tw.w.Write(chunk)
}
