package memory

import (
	"sync"

	"tsniper/internal/domain/scope"
	"tsniper/internal/domain/snipe"
)

var _ snipe.SnipeCache = (*SnipeCacheImpl)(nil)

type SnipeCacheImpl struct {
	mu             sync.RWMutex
	trackedIndices map[scope.AcademicScope][]string
	trackedCount   map[scope.AcademicScope]map[string]int
}

func NewSnipeCache() *SnipeCacheImpl {
	return &SnipeCacheImpl{
		trackedIndices: make(map[scope.AcademicScope][]string),
		trackedCount:   make(map[scope.AcademicScope]map[string]int),
	}
}

func (c *SnipeCacheImpl) Add(snp *snipe.Snipe) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.trackedCount[snp.Scope] == nil {
		c.trackedCount[snp.Scope] = make(map[string]int)
	}
	c.trackedCount[snp.Scope][snp.Index] += 1
	if c.trackedCount[snp.Scope][snp.Index] == 1 {
		c.trackedIndices[snp.Scope] = append(c.trackedIndices[snp.Scope], snp.Index)
	}
}

func (c *SnipeCacheImpl) Remove(snp *snipe.Snipe) {
	c.mu.Lock()
	defer c.mu.Unlock()
	counts, ok := c.trackedCount[snp.Scope]
	if !ok {
		return
	}
	counts[snp.Index] -= 1
	if counts[snp.Index] == 0 {
		delete(counts, snp.Index)
		indices := c.trackedIndices[snp.Scope]
		for i, v := range indices {
			if v == snp.Index {
				indices[i] = indices[len(indices)-1]
				c.trackedIndices[snp.Scope] = indices[:len(indices)-1]
				break
			}
		}
		if len(counts) == 0 {
			delete(c.trackedCount, snp.Scope)
			delete(c.trackedIndices, snp.Scope)
		}
	}
}

func (c *SnipeCacheImpl) Clear(snipes []*snipe.Snipe) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, snp := range snipes {
		counts, ok := c.trackedCount[snp.Scope]
		if !ok {
			continue
		}
		counts[snp.Index] -= 1
		if counts[snp.Index] == 0 {
			delete(counts, snp.Index)
			indices := c.trackedIndices[snp.Scope]
			for i, v := range indices {
				if v == snp.Index {
					indices[i] = indices[len(indices)-1]
					c.trackedIndices[snp.Scope] = indices[:len(indices)-1]
					break
				}
			}
			if len(counts) == 0 {
				delete(c.trackedCount, snp.Scope)
				delete(c.trackedIndices, snp.Scope)
			}
		}
	}
}

func (c *SnipeCacheImpl) Tracked(s scope.AcademicScope) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	indices := c.trackedIndices[s]
	if indices == nil {
		return nil
	}
	tracked := make([]string, len(indices))
	copy(tracked, indices)
	return tracked
}
