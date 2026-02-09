package exception

var (
	ErrRuyiExpandFailed  = Errorf("ruyi expand failed")
	ErrRuyiIsBigEnough   = Wrapf(ErrRuyiExpandFailed, "ruyi is big enough")
	ErrRuyiShrinkFailed  = Errorf("ruyi shrink failed")
	ErrRuyiIsSmallEnough = Wrapf(ErrRuyiShrinkFailed, "ruyi is small enough")
)

var (
	ErrInternal              = Errorf("internal error")
	ErrNoSupportedConverter  = Errorf("no supported converter")
	ErrConvertFailed         = Errorf("convert failed")
	ErrIllegalConverterParam = Errorf("illegal converter param") // 非法的转换器参数
)
