package attributes

type Attr struct {
	Depth, Retry int
}

type Storage map[int]*Attr

func (s Storage) Retry(id int) int {
	if attr, ok := s[id]; ok {
		return attr.Retry
	}

	return -1
}

func (s Storage) Depth(id int) int {
	if attr, ok := s[id]; ok {
		return attr.Depth
	}

	return -1
}

func (s Storage) Insert(id int, attr *Attr) {
	s[id] = attr
}

func (s Storage) Delete(id int) {
	delete(s, id)
}

func New() Storage {
	return map[int]*Attr{}
}
