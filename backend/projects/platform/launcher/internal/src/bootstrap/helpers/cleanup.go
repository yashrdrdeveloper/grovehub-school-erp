package helpers

import "sync"

// Cleanup releases owned resources in reverse initialization order.
type Cleanup struct { callbacks []func(); once sync.Once }
func(c *Cleanup) Add(fn func()){if fn!=nil {c.callbacks=append(c.callbacks,fn)}}
func(c *Cleanup) Close(){c.once.Do(func(){for i:=len(c.callbacks)-1;i>=0;i--{c.callbacks[i]()}})}
