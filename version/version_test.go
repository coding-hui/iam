package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsOfficialKubeVelaVersion(t *testing.T) {
	assert.Equal(t, true, IsOfficialIAMVersion("v1.0.3"))
	assert.Equal(t, true, IsOfficialIAMVersion("1.0.3"))
	assert.Equal(t, true, IsOfficialIAMVersion("v1.0"))
	assert.Equal(t, true, IsOfficialIAMVersion("v1.0+oauth2"))
	assert.Equal(t, false, IsOfficialIAMVersion("v1.-"))
}

func TestGetVersion(t *testing.T) {
	version, err := GetOfficialIAMVersion("v1.0.90")
	assert.Nil(t, err)
	assert.Equal(t, "1.0.90", version)
	version, err = GetOfficialIAMVersion("1.0.90")
	assert.Nil(t, err)
	assert.Equal(t, "1.0.90", version)
	version, err = GetOfficialIAMVersion("v1.0+oauth2")
	assert.Nil(t, err)
	assert.Equal(t, "1.0.0", version)
}
