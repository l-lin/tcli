package trello

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

func sanitize(name string) string {
	return strings.ReplaceAll(name, " ", "\\ ")
}

// getPos convert the given pos into appropriate type supported by Trello: either a string or a float
func getPos(in string) interface{} {
	if in == "top" || in == "bottom" {
		return in
	}
	pos, err := strconv.ParseFloat(in, 64)
	if err != nil {
		log.Debug().
			Str("pos", in).
			Err(err).
			Msg("could not parse pos to float")
		return in
	}
	return pos
}

// toTCliID converts a Trello entity into a unique ID understandable by tcli
// it's using the name, for the user experience in the completion, and the short link
// instead of the id to prevent having long lines in the completion.
func toTCliID(name, uid string) string {
	return fmt.Sprintf("%s[%s]", name, uid)
}
