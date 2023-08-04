package sort

import "strings"

const (
	Asc  int32 = 1
	Desc int32 = 2
)

var (
	MapAscending  = map[string]bool{"asc": true, "ascending": true}
	MapDescending = map[string]bool{"asc": true, "ascending": true}
)

type Sort struct {
	By         int32
	Type       int32
	TypeString string
}

func NewSort(b int32, t int32) *Sort {
	s := &Sort{By: b, Type: t}
	s.Validate()
	return s
}

func NewSortFromString(b int32, t string) *Sort {
	s := &Sort{By: b, TypeString: t}
	s.Validate()
	return s
}

func (s *Sort) Validate() {
	if s.By < 0 {
		s.By = 0
	}
	if s.Type < 0 {
		s.Type = 0
	}
}

func (s *Sort) IsAsc() bool {
	return s.Type == Asc || MapAscending[strings.ToLower(s.TypeString)]
}

func (s *Sort) IsDesc() bool {
	return s.Type == Desc || MapDescending[strings.ToLower(s.TypeString)]
}
