package filetransformer

type TransformType int

const (
	TransformTypeUnknown  = 0
	TransformTypeBasic    = 1
	TransformTypeLuno     = 2
	transformTypeSentinel = 3
)

func ValidTransformTypes() []TransformType {
	return []TransformType{
		TransformTypeUnknown,
		TransformTypeBasic,
		TransformTypeLuno,
	}
}

var transformTypeStrings = map[TransformType]string{
	TransformTypeUnknown: "Unknown",
	TransformTypeBasic:   "Basic (Don't use)",
	TransformTypeLuno:    "Luno",
}

func (tt TransformType) String() string {
	str, ok := transformTypeStrings[tt]
	if ok {
		return str
	}

	return "Unknown"
}

func (tt TransformType) Int() int {
	return int(tt)
}
