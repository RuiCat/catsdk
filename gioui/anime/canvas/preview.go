package canvas

import (
	"fmt"
	"os"
)

func loadFont(name string, style FontStyle) ([]byte, error) {
	filename, ok := FindSystemFont(name, style)
	if !ok {
		return nil, fmt.Errorf("failed to find font '%s'", name)
	}
	return os.ReadFile(filename)
}
