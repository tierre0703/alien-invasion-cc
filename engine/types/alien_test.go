package types
import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewAlien(t *testing.T) {
	a := NewAlien(100)
	require.Equal(t, 100, a.AlienID)
	require.Equal(t, "Alien #100", a.String())
}