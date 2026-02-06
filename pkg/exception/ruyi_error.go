package exception

var (
	ErrRuyiExpandFailed  = Errorf("ruyi expand failed")
	ErrRuyiIsBigEnough   = Wrapf(ErrRuyiExpandFailed, "ruyi is big enough")
	ErrRuyiShrinkFailed  = Errorf("ruyi shrink failed")
	ErrRuyiIsSmallEnough = Wrapf(ErrRuyiShrinkFailed, "ruyi is small enough")
)

var (
	ErrInternal               = Errorf("internal error")
	ErrNoSupportedConverter   = Errorf("no supported converter")
	ErrConvertFailed          = Errorf("convert failed")
	ErrInvalidConverterOutput = Errorf("converter returned invalid output type")
)
