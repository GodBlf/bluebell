package snowflake

import (
	"bluebell/initialize"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnowflake(t *testing.T) {
	initialize.Initialize()
	snowflake, err := NewSnowflake("2024-01-01", 1)
	assert.NotNil(t, snowflake, "Snowflake instance should not be nil")
	assert.NoError(t, err, "Failed to create Snowflake instance")
	id := snowflake.GetID()
	assert.NotZero(t, id, "Generated ID should not be zero")
}
