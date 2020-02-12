package parser


const (

    RootSpanConst = "null" // when we get a call from the outside (“null”)
	dateRegexp     = `[\dTZ\.\-\:]*`
	traceRegExp    = `[\w-]{1,}`
    serviceRegExp  = `[\w-]{1,}`
    spanRegExp     = `[\w-]{1,}`
)
