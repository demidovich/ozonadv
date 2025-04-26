package storage

import (
	"os"
	"testing"

	"github.com/demidovich/ozonadv/test/factory"
	"github.com/stretchr/testify/assert"
)

func TestCabinets_Add(t *testing.T) {
	storage := newTestingCabinetsStorage()
	cabinet := factory.Cabinet().New()

	storage.Add(cabinet)

	file := testingCabinetsFile()
	assert.FileExists(t, file)
}

func TestCabinets_All(t *testing.T) {
	storage := newTestingCabinetsStorage()
	cabinet := factory.Cabinet().New()
	storage.Add(cabinet)

	all := storage.All()

	assert.Len(t, all, 1)
	assert.Equal(t, cabinet.UUID, all[0].UUID)
}

func TestCabinets_Has(t *testing.T) {
	storage := newTestingCabinetsStorage()
	cabinet1 := factory.Cabinet().New()
	cabinet2 := factory.Cabinet().New()

	storage.Add(cabinet1)

	assert.True(t, storage.Has(cabinet1))
	assert.False(t, storage.Has(cabinet2))
}

func TestCabinets_Remove(t *testing.T) {
	storage := newTestingCabinetsStorage()
	cabinet := factory.Cabinet().New()

	storage.Add(cabinet)
	assert.Len(t, storage.All(), 1)

	storage.Remove(cabinet)
	assert.Empty(t, storage.All())
}

func newTestingCabinetsStorage() *cabinetsStorage {
	file := testingCabinetsFile()
	os.Remove(file)

	return newCabinetsStorage(file)
}

func testingCabinetsFile() string {
	return os.TempDir() + "/test-ozonadv-cabinets.json"
}
