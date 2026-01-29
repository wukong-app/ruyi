package exception

var (
	ErrRuyiExpandFailed  = New("ruyi expand failed")
	ErrRuyiIsBigEnough   = Wrap(ErrRuyiExpandFailed, "ruyi is big enough")
	ErrRuyiShrinkFailed  = New("ruyi shrink failed")
	ErrRuyiIsSmallEnough = Wrap(ErrRuyiShrinkFailed, "ruyi is small enough")
)
