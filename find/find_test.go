package find

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeName(t *testing.T) {
	assert.Equal(t, "/tmp/Artist/Album", sanitizePath("/tmp/Artist/Album/CD1/test.mp3"), "they should be equal")
	assert.Equal(t, "/tmp/Artist/Album", sanitizePath("/tmp/Artist/Album/CD10/test.mp3"), "they should be equal")
}

func TestValidFile(t *testing.T) {
	assert.Equal(t, true, isValidFile("/tmp/Artist/Album/CD1/test.mp3"), "they should be valid")
	assert.Equal(t, true, isValidFile("/tmp/Artist/Album/CD1/test.MP3"), "they should be valid")
	assert.Equal(t, true, isValidFile("/tmp/Artist/Album/CD1/test.alac"), "they should be valid")
	assert.Equal(t, true, isValidFile("/tmp/Artist/Album/CD1/test.flac"), "they should be valid")
	assert.Equal(t, true, isValidFile("/tmp/Artist/Album/CD1/test.m4p"), "they should be valid")
	assert.Equal(t, false, isValidFile("/tmp/Artist/Album/CD1/test.jpg"), "they should not be valid")
}
