package bufferpool

import "github.com/ruinwCN/morm/buffer"

var (
	_pool = buffer.NewPool()
	// Get retrieves a buffer from the pool, creating one if necessary.
	Get = _pool.Get
)
