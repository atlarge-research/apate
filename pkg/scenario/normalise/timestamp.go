package normalise

import (
	"time"
)

// Takes a timestamp as given in a public.scenario and desugars it into an integer milliseconds.
// Inputs can be strings containing
// * Integers 					Time is measured in seconds after the start of the scenario
// * Integers with postfixes (not case sensitive):
//		- <int/float>s			Seconds after the start of the scenario. Can be a float but never more precise than ms.
//		- <int>ms				Milliseconds after the start of the scenario
// 		- <int>m				Minutes after the start of the scenario
//		- <int>h				Hours after the start of the scenario
func desugarTimestamp(t string) (int, error) {
	duration, err := time.ParseDuration(t)
	if err != nil {
		return 0, err
	}

	return int(duration.Milliseconds()), nil
}
