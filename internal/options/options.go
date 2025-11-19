package options

type Options struct {
	Format         string
	Input          string
	Output         string
	TypeName       string
	PackageName    string
	UseStdin       bool
	Tags           []string
	OmitEmpty      bool
	SkipUnparsable bool
}

func NewOptions() *Options {
	return &Options{
		TypeName:       "Root",
		PackageName:    "main",
		Tags:           []string{},
		OmitEmpty:      false,
		SkipUnparsable: true,
	}
}
