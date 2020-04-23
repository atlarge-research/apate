package normalize

import (
	"fmt"
	"strconv"
	"strings"
)

// Takes a timestamp as given in a public.scenario and desugars it into an integer milliseconds.
// Inputs can be strings containing
// * Integers 					Time is measured in seconds after the start of the scenario
// * Integers with postfixes (not case sensitive):
//		- <int/float>s			Seconds after the start of the scenario. Can be a float but never more precise than ms.
//		- <int>ms				Milliseconds after the start of the scenario
// 		- <int>m				Minutes after the start of the scenario
//		- <int>h				Hours after the start of the scenario
func desugarTimestamp(time string) (int, error) {

	// First remove all leading and trailing spaces
	time = strings.TrimSpace(time)

	// postfix 1 and 2
	p1 := time[len(time) - 1]
	p2 := time[len(time) - 2]

	var v int

	switch p1 {
	case 's':
		switch p2 {
		case 'm':
			// Milliseconds (take everything before the second to last char and convert to int)

			iv, err := strconv.ParseInt(time[:len(time)-2], 10, 64)
			if err != nil {
				return 0, err
			}

			v = int(iv)

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
			// I match on numbers or . here (. for floats) to be more descriptive in error messages.
			// Otherwise invalid postfixes would result in integer conversion errors.
			// Seconds (take everything before the last char and convert to float)

			fv, err := strconv.ParseFloat(time[:len(time)-1], 64)
			if err != nil {
				return 0, err
			}

			// multiply by 1000 to go to milliseconds, then round to the nearest ms
			v = int(fv * 1000)

		default:
			return 0, fmt.Errorf("Couldn't decode postfix. Possible values are ms, s, m, h or no postfix. (%s)", time)
		}
	case 'm':
		// Minutes (take everything before the last char and convert to int)

		iv, err := strconv.ParseInt(time[:len(time)-1], 10, 64)
		if err != nil {
			return 0, err
		}

		// multiply by 1000 * 60 to go to milliseconds.
		v = int(iv * 60 * 1000)
	case 'h':
		// Minutes (take everything before the last char and convert to int)

		iv, err := strconv.ParseInt(time[:len(time)-1], 10, 64)
		if err != nil {
			return 0, err
		}

		// multiply by 1000 * 60 * 60 to go to milliseconds.
		v = int(iv * 60  * 60 * 1000)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
		// I match on numbers or . here (. for floats) to be more descriptive in error messages.
		// Otherwise invalid postfixes would result in integer conversion errors.
		// Seconds (no postfix, take everything and convert to float)

		fv, err := strconv.ParseFloat(time, 64)
		if err != nil {
			return 0, err
		}

		// multiply by 1000 to go to milliseconds, then round to the nearest ms
		v = int(fv * 1000)
	default:
		return 0, fmt.Errorf("Couldn't decode postfix. Possible values are ms, s, m, h or no postfix. (%s)", time)
	}

	return v, nil
}