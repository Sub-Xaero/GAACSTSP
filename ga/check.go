package ga

import "log"

// Standard error checking function, used to crash the program if any errors encountered
// works for all errors
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
