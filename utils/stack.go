package utils

type Stack struct {
	data map[int]interface{}
}

func New() *Stack {
	s := new(Stack)
	s.data = make(map[int]interface{})
	return s
}

func (s *Stack) Push(data interface{}) {
	s.data[len(s.data)] = data
}

func (s *Stack) Pop() {
	delete(s.data, len(s.data)-1)
}

func (s *Stack) String() string {
	info := ""
	for i := 0; i < len(s.data); i++ {
		info = info + "[" + AnyToString(s.data[i]) + "]"
	}
	return info
}
