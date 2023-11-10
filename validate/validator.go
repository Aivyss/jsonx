package validate

type Validator[V any] interface {
	Validate(v V) error
}

type OrderedValidator[V any] interface {
	Validator[V]
	Order() int
}
