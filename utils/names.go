package utils

import (
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	mrand "math/rand"
	"regexp"
	"unicode/utf8"
)

var firstNames = []string{
	"john", "james", "robert", "michael", "william",
	"david", "richard", "joseph", "thomas", "charles",
	"daniel", "matthew", "anthony", "mark", "steven",
	"paul", "andrew", "joshua", "kevin", "brian",
	"george", "edward", "timothy", "jason", "ryan",
	"jacob", "nicholas", "eric", "jonathan", "justin",
	"scott", "brandon", "benjamin", "samuel", "patrick",
	"jack", "tyler", "aaron", "henry", "adam",
	"nathan", "kyle", "jeremy", "sean", "ethan",
	"noah", "jordan", "dylan", "gabriel", "vincent",

	"oliver", "leo", "liam", "lucas", "felix",
	"max", "emil", "anton", "leon", "lukas",
	"tobias", "jonas", "simon", "fabian", "marco",
	"sebastian", "manuel", "ivan", "nikola", "milan",
	"stefano", "lorenzo", "giovanni", "pierre", "luc",
	"antoine", "julien", "carlos", "miguel", "diego",
	"javier", "antonio", "rafael", "pablo", "andres",
	"oleg", "dmitri", "alexei", "roman", "kasper",
	"anders", "lars", "mikkel", "henrik", "oskar", "christian",

	"emma", "olivia", "ava", "sophia", "isabella",
	"mia", "amelia", "charlotte", "harper", "evelyn",
	"abigail", "emily", "elizabeth", "sofia", "avery",
	"scarlett", "grace", "chloe", "victoria", "riley",
	"arabella", "lily", "hannah", "ella", "nora",
	"zoe", "lila", "clara", "julia", "sarah",
	"maria", "anna", "kate", "paula", "laura",
	"lucia", "isabel", "camila", "alice", "amelie",
	"sophie", "leonie", "marie", "mia", "emilia",
	"eva", "elena", "katharina", "lena", "julia",
}

var adjectives = []string{
	"brave", "calm", "eager", "fancy", "gentle",
	"happy", "jolly", "kind", "lucky", "mighty",
	"nice", "proud", "quick", "sharp", "smart",
	"sunny", "swift", "wise", "bold", "cool",
}

// generates a random hex suffix of given byte length
func randomSuffix(bytes int) string {
	b := make([]byte, bytes)
	_, _ = crand.Read(b) // crypto/rand
	return hex.EncodeToString(b)
}

// generate a pull secret name like "pullsecret-brave-john-1a2b"
func GeneratePullSecretName() string {
	adjective := adjectives[mrand.Intn(len(adjectives))]
	name := firstNames[mrand.Intn(len(firstNames))]
	suffix := randomSuffix(2) // 4 hex chars

	return fmt.Sprintf(
		"pullsecret-%s-%s-%s",
		adjective,
		name,
		suffix,
	)
}

func IsK8sNameValid(name string) bool {
	if utf8.RuneCountInString(name) > 253 {
		return false
	}

	// Regex: start/end alphanumeric, only a-z, 0-9, -
	k8sNameRegex := `^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	matched, _ := regexp.MatchString(k8sNameRegex, name)
	return matched
}
