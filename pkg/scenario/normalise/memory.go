package normalise

import (
	"fmt"
	"strconv"
	"strings"
)

func desugarMemory(memory string) (int64, error) {
	// First remove all leading and trailing spaces
	memory = strings.TrimSpace(memory)
	// Then add exactly 3 spaces only at the front. This  will make it so each string slice operation succeeds.
	memory = "   " + memory

	var v int64

	em := fmt.Errorf("could not parse memory size %s", memory)

	switch strings.ToUpper(memory[len(memory)-3:]) {
	case "GIB":
		iv, err := strconv.ParseInt(memory[3:len(memory)-3], 10, 64)
		if err != nil {
			return 0, em
		}
		v = iv * 1024 * 1024 * 1024
	case "MIB":
		iv, err := strconv.ParseInt(memory[3:len(memory)-3], 10, 64)
		if err != nil {
			return 0, em
		}
		v = iv * 1024 * 1024
	case "KIB":
		iv, err := strconv.ParseInt(memory[3:len(memory)-3], 10, 64)
		if err != nil {
			return 0, em
		}
		v = iv * 1024
	default:
		switch strings.ToUpper(memory[len(memory)-2:]) {
		case "GB":
			iv, err := strconv.ParseInt(memory[3:len(memory)-2], 10, 64)
			if err != nil {
				return 0, em
			}
			v = iv * 1000 * 1000 * 1000
		case "MB":
			iv, err := strconv.ParseInt(memory[3:len(memory)-2], 10, 64)
			if err != nil {
				return 0, em
			}
			v = iv * 1000 * 1000
		case "KB":
			iv, err := strconv.ParseInt(memory[3:len(memory)-2], 10, 64)
			if err != nil {
				return 0, em
			}
			v = iv * 1000
		default:
			switch strings.ToUpper(memory[len(memory)-1:]) {
			case "G":
				iv, err := strconv.ParseInt(memory[3:len(memory)-1], 10, 64)
				if err != nil {
					return 0, em
				}
				v = iv * 1000 * 1000 * 1000
			case "M":
				iv, err := strconv.ParseInt(memory[3:len(memory)-1], 10, 64)
				if err != nil {
					return 0, em
				}
				v = iv * 1000 * 1000
			case "K":
				iv, err := strconv.ParseInt(memory[3:len(memory)-1], 10, 64)
				if err != nil {
					return 0, em
				}
				v = iv * 1000
			case "B":
				iv, err := strconv.ParseInt(memory[3:len(memory)-1], 10, 64)
				if err != nil {
					return 0, em
				}
				v = iv
			default:
				iv, err := strconv.ParseInt(memory[3:], 10, 64)
				if err != nil {
					return 0, em
				}
				v = iv
			}
		}
	}

	return v, nil
}
