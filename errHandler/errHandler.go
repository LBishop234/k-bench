package errhandler

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Takes a error and a clean error message.
// Logs the error and returns a combined error message.
// This exists to reduce having to rewrite this all over the place.
func Error(cleanErr string, errors ...error) error {
	if len(errors) == 0 && cleanErr == "" {
		log.Warn("Error() function recieved blank  cleanErr value and no errors")
		return nil
	} else if len(errors) == 0 {
		err := fmt.Errorf(cleanErr)
		log.Error(err)
		return err
	} else if cleanErr == "" {
		err := concatErrors(errors)
		log.Error(err)
		return err
	} else {
		err := fmt.Errorf("%s | %s", cleanErr, concatErrors(errors))
		log.Error(err)
		return err
	}
}

// Combines an array of errors into a single error
func concatErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	} else if len(errs) == 1 {
		return errs[0]
	} else {
		err := errs[0]

		for i := 1; i < len(errs); i++ {
			err = fmt.Errorf("%s & %s", err, errs[0])
		}

		return err
	}
}
