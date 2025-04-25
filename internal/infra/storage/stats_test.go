package storage

import (
	"fmt"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/demidovich/ozonadv/pkg/utils"
	"github.com/demidovich/ozonadv/test/factory"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	storage := newTestingStatsStorage()
	stat := factory.Stat().New()

	storage.Add(stat)

	file := fmt.Sprintf("%s/%s.json", storage.dir, stat.UUID)
	assert.FileExists(t, file)
}

func TestRemove(t *testing.T) {
	storage := newTestingStatsStorage()
	stat := factory.Stat().New()

	storage.Add(stat)
	storage.Remove(stat)

	file := fmt.Sprintf("%s/%s.json", storage.dir, stat.UUID)
	assert.NoFileExists(t, file)
}

func TestAll(t *testing.T) {
	storage := newTestingStatsStorage()
	stat := factory.Stat().New()

	storage.Add(stat)
	all := storage.All()

	assert.Len(t, all, 1)
	assert.Equal(t, stat.UUID, all[0].UUID)
}

func TestAddDownloadsFile(t *testing.T) {
	storage := newTestingStatsStorage()
	stat := factory.Stat().New()
	fname := gofakeit.UUID() + ".csv"
	fdata := gofakeit.UUID()

	storage.AddDownloadsFile(stat, fname, []byte(fdata))

	file := fmt.Sprintf("%s/downloads/%s", storage.dir, fname)
	assert.FileExists(t, file)

	data, err := os.ReadFile(file)
	if assert.NoError(t, err) {
		assert.Equal(t, fdata, string(data))
	}
}

func TestReadDownloadsFile(t *testing.T) {
	storage := newTestingStatsStorage()
	stat := factory.Stat().New()
	fname := gofakeit.UUID() + ".csv"
	fdata := gofakeit.UUID()

	storage.AddDownloadsFile(stat, fname, []byte(fdata))
	data := storage.ReadDownloadedFile(stat, fname)

	assert.Equal(t, fdata, string(data))
}

func newTestingStatsStorage() *statsStorage {
	dir := "/tmp/test-ozonadv-stats"
	os.RemoveAll(dir)
	utils.DirInitOrFail(dir)

	return newStatsStorage(dir)
}
