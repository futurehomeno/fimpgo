package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEcdsaKey_JsonRoundTrip(t *testing.T) {
	orig := NewEcdsaKey()
	require.NoError(t, orig.Generate())
	priv, pub, err := orig.ExportJsonEncodedKeys()
	require.NoError(t, err)

	signer := NewEcdsaKey()
	require.NoError(t, signer.ImportJsonPrivateKey(priv))
	verifier := NewEcdsaKey()
	require.NoError(t, verifier.ImportJsonPublicKey(pub))

	assert.True(t, orig.PrivateKey().Equal(signer.PrivateKey()))
	assert.True(t, orig.PublicKey().Equal(verifier.PublicKey()))
}

func TestEcdsaKey_ImportJsonPublicKey_RejectsOffCurvePoint(t *testing.T) {
	err := NewEcdsaKey().ImportJsonPublicKey(JsonEcKey{X: "01", Y: "01"})
	assert.Error(t, err)
}
