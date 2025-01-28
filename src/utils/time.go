package utils

import "time"

/*
*
* Returns the current date and time
@returns string
*/
func GetCurrentTime() string {
	return time.Now().Format(time.RFC3339)
}
