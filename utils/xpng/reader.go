package xpng

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"slices"
)

type ImageMetadata struct {
	Width  uint32
	Height uint32
	Txt    map[string]string
}

const maxTxtSize = 1 << 16
const pngHeader = "\x89PNG\r\n\x1a\n"

func DecodeTxt(r io.Reader) (ImageMetadata, error) {
	var buff [maxTxtSize]byte

	result := ImageMetadata{Txt: make(map[string]string)}

	if _, err := io.ReadFull(r, buff[:len(pngHeader)]); err != nil {
		return ImageMetadata{}, err
	}

	if string(buff[:len(pngHeader)]) != pngHeader {
		return ImageMetadata{}, errors.New("not a PNG file")
	}

	for {
		if _, err := io.ReadFull(r, buff[:8]); err != nil {
			return ImageMetadata{}, err
		}

		length := binary.BigEndian.Uint32(buff[:4])

		switch string(buff[4:8]) {
		case "IHDR":
			if length != 13 {
				return ImageMetadata{}, errors.New("bad IHDR length")
			}

			if _, err := io.ReadFull(r, buff[:length]); err != nil {
				return ImageMetadata{}, err
			}

			result.Width = binary.BigEndian.Uint32(buff[0:4])
			result.Height = binary.BigEndian.Uint32(buff[4:8])
		case "iTXt":
			if length > maxTxtSize {
				return ImageMetadata{}, errors.New("iTXt chunk is too large")
			}

			if _, err := io.ReadFull(r, buff[:length]); err != nil {
				return ImageMetadata{}, err
			}

			sepIndex := slices.Index(buff[:length], 0)
			if sepIndex < 0 {
				break
			}

			key := string(buff[:sepIndex])
			value := string(buff[sepIndex+5 : length])
			result.Txt[key] = value
		case "IDAT":
			return result, nil
		case "IEND":
			return result, nil
		default:
		}

		if _, err := io.ReadFull(r, buff[:4]); err != nil {
			return ImageMetadata{}, err
		}
	}
}

func FileMetadata(fileName string) (ImageMetadata, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return ImageMetadata{}, err
	}

	defer file.Close()

	metadata, err := DecodeTxt(file)
	if err != nil {
		return ImageMetadata{}, err
	}

	return metadata, nil
}
