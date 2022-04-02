//go:build !prod

package force

import "github.com/starboard-nz/go-rebuild"

func init() {
	rebuild.Run(true)
}
