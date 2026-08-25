package render

import (
	"fmt"
	"sync"
)

const RenderVersion = 1 // Bump this when HTML/CSS templates change

type FrameCache struct {
	mu     sync.RWMutex
	frames map[string][][]byte // key: "chapter:{id}:v{version}"
}

func NewFrameCache() *FrameCache {
	return &FrameCache{
		frames: make(map[string][][]byte),
	}
}

func (c *FrameCache) key(chapterID int) string {
	return fmt.Sprintf("chapter:%d:v%d", chapterID, RenderVersion)
}

func (c *FrameCache) Get(chapterID int, subPageIndex int) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	chapterFrames, exists := c.frames[c.key(chapterID)]
	if !exists || subPageIndex < 0 || subPageIndex >= len(chapterFrames) {
		return nil, false
	}
	return chapterFrames[subPageIndex], true
}

func (c *FrameCache) Set(chapterID int, frames [][]byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.frames[c.key(chapterID)] = frames
}

func (c *FrameCache) Invalidate(chapterID int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.frames, c.key(chapterID))
}

func (c *FrameCache) FrameCount(chapterID int) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.frames[c.key(chapterID)])
}

func (c *FrameCache) GetAllFrames(chapterID int) ([][]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	frames, exists := c.frames[c.key(chapterID)]
	if !exists || len(frames) == 0 {
		return nil, false
	}
	return frames, true
}
