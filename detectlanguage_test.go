package detectlanguage_test

import (
	"os"

	"github.com/4O4-Not-F0und/detectlanguage-go"
)

var client = detectlanguage.New(os.Getenv("DETECTLANGUAGE_API_KEY"))
