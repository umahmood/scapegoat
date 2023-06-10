package scapegoat_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/umahmood/scapegoat"
)

func TestWithFloats(t *testing.T) {
	assert := assert.New(t)

	sg, err := scapegoat.New[float64](scapegoat.DefaultAlpha)
	assert.NotNil(sg)
	assert.Nil(err)

	keys := []float64{42.3, 99.1, 0.3}
	for _, key := range keys {
		err := sg.Insert(key)
		assert.Nil(err)
	}
}

func TestInsertDuplicateKey(t *testing.T) {
	assert := assert.New(t)

	sg, err := scapegoat.New[int](scapegoat.DefaultAlpha)
	assert.NotNil(sg)
	assert.Nil(err)

	keys := []int{42, 99, 3}
	for _, key := range keys {
		err = sg.Insert(key)
		assert.Nil(err)
	}

	err = sg.Insert(99)
	assert.Nil(err)

	assert.Equal(uint64(0), sg.Stats.TotalRebalances)
	assert.Equal(uint64(0), sg.Stats.TotalRebalancesAfterInsert)
	assert.Equal(uint64(0), sg.Stats.TotalRebalancesAfterRemove)
	assert.Equal(uint64(len(keys)), sg.Stats.TotalInserts)
	assert.Equal(uint64(0), sg.Stats.TotalRemovals)
	assert.Equal(uint64(0), sg.Stats.TotalSearches)
}

func TestInsertIntoEmptyTree(t *testing.T) {
	assert := assert.New(t)

	sg, err := scapegoat.New[int](scapegoat.DefaultAlpha)

	assert.NotNil(sg)
	assert.Nil(err)

	err = sg.Insert(42)
	assert.Nil(err)

	assert.Equal(uint64(0), sg.Stats.TotalRebalances)
	assert.Equal(uint64(0), sg.Stats.TotalRebalancesAfterInsert)
	assert.Equal(uint64(0), sg.Stats.TotalRebalancesAfterRemove)
	assert.Equal(uint64(1), sg.Stats.TotalInserts)
	assert.Equal(uint64(0), sg.Stats.TotalRemovals)
	assert.Equal(uint64(0), sg.Stats.TotalSearches)
}

func TestInsertNoRebalance(t *testing.T) {
	assert := assert.New(t)

	sg, err := scapegoat.New[int](scapegoat.DefaultAlpha)

	assert.NotNil(sg)
	assert.Nil(err)

	keys := []int{42, 99, 3}
	for _, key := range keys {
		err := sg.Insert(key)
		assert.Nil(err)
	}
	assert.Equal(uint64(0), sg.Stats.TotalRebalances)
	assert.Equal(uint64(0), sg.Stats.TotalRebalancesAfterInsert)
	assert.Equal(uint64(0), sg.Stats.TotalRebalancesAfterRemove)
	assert.Equal(uint64(len(keys)), sg.Stats.TotalInserts)
	assert.Equal(uint64(0), sg.Stats.TotalRemovals)
	assert.Equal(uint64(0), sg.Stats.TotalSearches)
}

func TestInsertTriggersRebalance(t *testing.T) {
	assert := assert.New(t)

	sg, err := scapegoat.New[int](scapegoat.DefaultAlpha)

	assert.NotNil(sg)
	assert.Nil(err)

	keys := []int{42, 27, 3, 1, 2, 99, 55, 48, 39, 47, 46}
	for _, key := range keys {
		err := sg.Insert(key)
		assert.Nil(err)
	}
	assert.Equal(uint64(1), sg.Stats.TotalRebalances)
	assert.Equal(uint64(1), sg.Stats.TotalRebalancesAfterInsert)
	assert.Equal(uint64(0), sg.Stats.TotalRebalancesAfterRemove)
	assert.Equal(uint64(len(keys)), sg.Stats.TotalInserts)
	assert.Equal(uint64(0), sg.Stats.TotalRemovals)
	assert.Equal(uint64(0), sg.Stats.TotalSearches)
}

func TestRemoveNoRebalance(t *testing.T) {
	assert := assert.New(t)

	sg, err := scapegoat.New[int](scapegoat.DefaultAlpha)
	assert.NotNil(sg)
	assert.Nil(err)

	keys := []int{42, 99, 3}
	for _, key := range keys {
		err := sg.Insert(key)
		assert.Nil(err)
	}
	wasRemoved := sg.Remove(99)
	assert.True(wasRemoved)

	assert.Equal(uint64(0), sg.Stats.TotalRebalances)
	assert.Equal(uint64(0), sg.Stats.TotalRebalancesAfterInsert)
	assert.Equal(uint64(0), sg.Stats.TotalRebalancesAfterRemove)
	assert.Equal(uint64(len(keys)), sg.Stats.TotalInserts)
	assert.Equal(uint64(1), sg.Stats.TotalRemovals)
	assert.Equal(uint64(0), sg.Stats.TotalSearches)
}

func TestRemoveTriggersRebalance(t *testing.T) {
	assert := assert.New(t)

	sg, err := scapegoat.New[int](scapegoat.DefaultAlpha)
	assert.NotNil(sg)
	assert.Nil(err)

	keys := []int{42, 27, 3, 1, 2, 99, 55, 48, 39, 47, 46, 45}
	for _, key := range keys {
		err := sg.Insert(key)
		assert.Nil(err)
	}
	keysToRemove := []int{3, 1, 2, 99, 55, 48, 39, 47}
	for _, key := range keysToRemove {
		wasRemoved := sg.Remove(key)
		assert.True(wasRemoved)
	}
	assert.Equal(uint64(2), sg.Stats.TotalRebalances)
	assert.Equal(uint64(1), sg.Stats.TotalRebalancesAfterInsert)
	assert.Equal(uint64(1), sg.Stats.TotalRebalancesAfterRemove)
	assert.Equal(uint64(len(keys)), sg.Stats.TotalInserts)
	assert.Equal(uint64(len(keysToRemove)), sg.Stats.TotalRemovals)
	assert.Equal(uint64(0), sg.Stats.TotalSearches)
}

func TestAttemptToRemoveKeysThatDoesNotExist(t *testing.T) {
	assert := assert.New(t)

	sg, err := scapegoat.New[int](scapegoat.DefaultAlpha)
	assert.NotNil(sg)
	assert.Nil(err)

	keys := []int{42, 27, 3, 1, 2, 99, 55, 48, 39, 47, 46, 45}
	for _, key := range keys {
		err := sg.Insert(key)
		assert.Nil(err)
	}
	keysToRemove := []int{4, 50, 69, 256, 1337, 1888}
	for _, key := range keysToRemove {
		wasRemoved := sg.Remove(key)
		assert.False(wasRemoved)
	}
	assert.Equal(uint64(1), sg.Stats.TotalRebalances)
	assert.Equal(uint64(1), sg.Stats.TotalRebalancesAfterInsert)
	assert.Equal(uint64(0), sg.Stats.TotalRebalancesAfterRemove)
	assert.Equal(uint64(len(keys)), sg.Stats.TotalInserts)
	assert.Equal(uint64(0), sg.Stats.TotalRemovals)
	assert.Equal(uint64(0), sg.Stats.TotalSearches)
}

// TODO: need to ADD tests for remove cases.

func TestSearchKeyExists(t *testing.T) {
	assert := assert.New(t)

	sg, err := scapegoat.New[int](scapegoat.DefaultAlpha)
	assert.NotNil(sg)
	assert.Nil(err)

	keys := []int{42, 27, 3, 1, 2, 99, 55, 48, 39, 47, 46, 45}
	for _, key := range keys {
		err := sg.Insert(key)
		assert.Nil(err)
	}

	for _, key := range keys {
		exists := sg.Search(key)
		assert.True(exists)
	}
	assert.Equal(uint64(1), sg.Stats.TotalRebalances)
	assert.Equal(uint64(1), sg.Stats.TotalRebalancesAfterInsert)
	assert.Equal(uint64(0), sg.Stats.TotalRebalancesAfterRemove)
	assert.Equal(uint64(len(keys)), sg.Stats.TotalInserts)
	assert.Equal(uint64(0), sg.Stats.TotalRemovals)
	assert.Equal(uint64(len(keys)), sg.Stats.TotalSearches)
}

func TestSearchKeyDoesNotExist(t *testing.T) {
	assert := assert.New(t)

	sg, err := scapegoat.New[int](scapegoat.DefaultAlpha)
	assert.NotNil(sg)
	assert.Nil(err)

	keys := []int{42, 27, 3, 1, 2, 99, 55, 48, 39, 47, 46, 45}
	for _, key := range keys {
		err := sg.Insert(key)
		assert.Nil(err)
	}
	exists := sg.Search(256)
	assert.False(exists)

	assert.Equal(uint64(1), sg.Stats.TotalRebalances)
	assert.Equal(uint64(1), sg.Stats.TotalRebalancesAfterInsert)
	assert.Equal(uint64(0), sg.Stats.TotalRebalancesAfterRemove)
	assert.Equal(uint64(len(keys)), sg.Stats.TotalInserts)
	assert.Equal(uint64(0), sg.Stats.TotalRemovals)
	assert.Equal(uint64(1), sg.Stats.TotalSearches)
}

// TODO: test for lots of random inserts, removes, searches
func TestRandomInsertsRemovesAndSearches(t *testing.T) {
	randomKeys := func() []int {
		seed := time.Now().UnixNano()
		rand.New(rand.NewSource(seed))
		min := 0
		max := 1000
		count := 500
		keys := make([]int, count)
		for i := 0; i < count; i++ {
			keys[i] = rand.Intn(max-min+1) + min
		}
		return keys
	}

	assert := assert.New(t)

	sg, err := scapegoat.New[int](scapegoat.DefaultAlpha)
	assert.NotNil(sg)
	assert.Nil(err)

	var (
		totalRemovals uint64
		totalSearches uint64
		keys          = randomKeys()
	)
	// keys can contain duplicates which are not inserted. so,
	// totalInserts is not tracked.
	for _, key := range keys[:10] {
		err := sg.Insert(key)
		assert.Nil(err)
	}
	for _, key := range keys[10:] {
		r := rand.Intn(3)
		switch r {
		case 0:
			err := sg.Insert(key)
			assert.Nil(err)
		case 1:
			removed := sg.Remove(key)
			if removed {
				totalRemovals++
			}
		case 2:
			sg.Search(key)
			totalSearches++
		}
	}
	assert.GreaterOrEqual(sg.Stats.TotalRebalances, uint64(0))
	assert.GreaterOrEqual(sg.Stats.TotalRebalancesAfterInsert, uint64(0))
	assert.GreaterOrEqual(sg.Stats.TotalRebalancesAfterRemove, uint64(0))
	assert.Greater(sg.Stats.TotalInserts, uint64(1))
	assert.Equal(totalRemovals, sg.Stats.TotalRemovals)
	assert.Equal(totalSearches, sg.Stats.TotalSearches)
}

func TestAlphaIsNotLessThanOrEqualToZero(t *testing.T) {
	assert := assert.New(t)
	for _, alpha := range []float64{0, -0.1, -0.001, -1.5} {
		sg, err := scapegoat.New[int](alpha)
		assert.Nil(sg)
		assert.ErrorIs(err, scapegoat.AlphaValueErr)
	}
}

func TestValidAlphaRanges(t *testing.T) {
	assert := assert.New(t)

	randomAlphas := []float64{
		0.8493577575966338,
		0.30759034800057505,
		0.6632902098507345,
		1.5263767974394318,
		2.1326543695173865,
		2.150112146060614,
		3.1496852027022437,
		4.91098032940735,
		4.326886633697387,
		5.0,
	}
	for _, alpha := range randomAlphas {
		sg, err := scapegoat.New[int](alpha)
		assert.NotNil(sg)
		assert.NoError(err)
	}
}
