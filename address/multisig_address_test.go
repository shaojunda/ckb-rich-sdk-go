package address

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateSecp256k1MultisigScript(t *testing.T) {
	var publicKeys [][]byte

	key, err := hex.DecodeString("032edb83018b57ddeb9bcc7287c5cc5da57e6e0289d31c9e98cb361e88678d6288")
	if err != nil {
		assert.Error(t, err)
	}
	publicKeys = append(publicKeys, key)

	key, err = hex.DecodeString("033aeb3fdbfaac72e9e34c55884a401ee87115302c146dd9e314677d826375dc8f")
	if err != nil {
		assert.Error(t, err)
	}
	publicKeys = append(publicKeys, key)

	key, err = hex.DecodeString("029a685b8206550ea1b600e347f18fd6115bffe582089d3567bec7eba57d04df01")
	if err != nil {
		assert.Error(t, err)
	}
	publicKeys = append(publicKeys, key)

	script, err := GenerateSecp256k1MultisigScript(0, 2, publicKeys)
	if err != nil {
		assert.Error(t, err)
	}

	address, err := Generate(Testnet, script)
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, "ckt1qyqlqn8vsj7r0a5rvya76tey9jd2rdnca8lqh4kcuq", address)
}
