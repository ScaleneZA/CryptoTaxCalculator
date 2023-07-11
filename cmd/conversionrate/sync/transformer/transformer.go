package transformer

type Transformer interface {
	TransformRow(row string) error
}
