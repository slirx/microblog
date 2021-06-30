package template

type Mock struct {
	GenerateFn func(t GeneratorType, fileName string, data interface{}) (string, error)
}

func (m Mock) Generate(t GeneratorType, fileName string, data interface{}) (string, error) {
	return m.GenerateFn(t, fileName, data)
}
