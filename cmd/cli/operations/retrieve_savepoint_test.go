package operations

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

/*
 * RetrieveLatestSavepoint
 */
func TestRetrieveLatestSavepointShouldReturnAnErrorIfItCannotReadFromDir(t *testing.T) {
	operator := RealOperator{
		Filesystem: afero.NewMemMapFs(),
	}

	files, err := operator.retrieveLatestSavepoint("/savepoints")

	assert.Equal(t, "", files)
	assert.EqualError(t, err, "open /savepoints: file does not exist")
}

func TestRetrieveLatestSavepointShouldReturnAnTheNewestFile(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	filesystem.Mkdir("/savepoints/", 0755)
	afero.WriteFile(filesystem, "/savepoints/savepoint-683b3f-59401d30cfc4", []byte("file a"), 644)
	afero.WriteFile(filesystem, "/savepoints/savepoint-323b3f-59401d30eoe6", []byte("file b"), 644)

	operator := RealOperator{
		Filesystem: filesystem,
	}

	files, err := operator.retrieveLatestSavepoint("/savepoints")

	assert.Equal(t, "/savepoints/savepoint-323b3f-59401d30eoe6", files)
	assert.Nil(t, err)
}

func TestRetrieveLatestSavepointShouldRemoveTheTrailingSlashFromTheSavepointDirectory(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	filesystem.Mkdir("/savepoints/", 0755)
	afero.WriteFile(filesystem, "/savepoints/savepoint-683b3f-59401d30cfc4", []byte("file a"), 644)
	afero.WriteFile(filesystem, "/savepoints/savepoint-323b3f-59401d30eoe6", []byte("file b"), 644)

	operator := RealOperator{
		Filesystem: filesystem,
	}

	files, err := operator.retrieveLatestSavepoint("/savepoints/")

	assert.Equal(t, "/savepoints/savepoint-323b3f-59401d30eoe6", files)
	assert.Nil(t, err)
}

func TestRetrieveLatestSavepointShouldReturnAnErrorWhenDirEmpty(t *testing.T) {
	filesystem := afero.NewMemMapFs()
	filesystem.Mkdir("/savepoints/", 0755)

	operator := RealOperator{
		Filesystem: filesystem,
	}

	files, err := operator.retrieveLatestSavepoint("/savepoints")

	assert.Equal(t, "", files)
	assert.EqualError(t, err, "No savepoints present in directory: /savepoints")
}
