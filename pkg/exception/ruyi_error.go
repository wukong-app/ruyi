package exception

var (
	ErrRuyiExpandFailed  = Errorf("ruyi expand failed")
	ErrRuyiIsBigEnough   = Wrapf(ErrRuyiExpandFailed, "ruyi is big enough")
	ErrRuyiShrinkFailed  = Errorf("ruyi shrink failed")
	ErrRuyiIsSmallEnough = Wrapf(ErrRuyiShrinkFailed, "ruyi is small enough")
)
