package oxford

import "errors"

type SpeakError error

var (
	errUnknownLocale = SpeakError(errors.New("Unknown locale"))
)

func IsSpeakError(err error) bool {
	_, ok := err.(SpeakError)
	return ok
}
