package spinlock

import (
	"math/rand"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	iterationThreshold = 4 // the more this value the long we wait before begin running of "time.Sleep" in Lock()
)

const (
	// We leave positive values free for future needs (it's required to implement RLock in future; RLock will increment the value).
	unlocked = int32(0)
	locked   = -1
)

type Locker struct {
	state int32
}

func (l *Locker) Lock() {
	i := 0
	for !atomic.CompareAndSwapInt32(&l.state, unlocked, locked) {
		if i > iterationThreshold {
			time.Sleep(time.Nanosecond * time.Duration((50 + rand.Intn(950))))
		} else {
			runtime.Gosched()
		}
		i++
	}
}

func (l *Locker) Unlock() {
	if !atomic.CompareAndSwapInt32(&l.state, locked, unlocked) {
		panic(`Unlock()-ing non-locked locker`)
	}
}

func (l *Locker) SetUnlocked() {
	atomic.StoreInt32(&l.state, unlocked)
}

func init() {
	if runtime.NumCPU() == 1 {
		logrus.Warnf("Only one logical CPU. It's not recommended to use spinlocks here.")
		return
	}
}
