package ga

import "log"

// Standard error checking function, used across package to crash the program if any errors encountered. Results in a
// call to
// 	log.Fatal(err)
// In the instance that:
// 	err != nil
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
