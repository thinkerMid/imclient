package intervalTimer

// Timer 间隔器
type Timer struct {
	nowInterval int32
	second      int32
}

// New .
func New(second uint32) *Timer {
	return &Timer{second: int32(second)}
}

// ChangeWaitSecond .
func (t *Timer) ChangeWaitSecond(v uint32) {
	t.second = int32(v)
	t.Reset()
}

// Enabled .
func (t *Timer) Enabled() bool {
	return t.second > 0
}

// Timing .
func (t *Timer) Timing() {
	// 没设置时长不操作
	if t.second == 0 {
		return
	}

	t.nowInterval++
}

// TimingEnd .
func (t *Timer) TimingEnd() bool {
	// 没设置时长默认计时结束
	if t.second == 0 {
		return true
	}

	return t.nowInterval == t.second
}

// Reset .
func (t *Timer) Reset() {
	t.nowInterval = 0
}
