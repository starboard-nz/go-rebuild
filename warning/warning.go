//go:build !prod
// +build !prod

package warning

import "github.com/starboard-nz/go-rebuild"

func init() {
	rebuild.Run(false)
}
