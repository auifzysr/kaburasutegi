package infra

import "fmt"

type LocalRecord struct {
	records map[int]string
}

func (s *LocalRecord) Record(text string) error {
	return s.RecordAt(len(s.records), text)
}

func (s *LocalRecord) RecordAt(id int, text string) error {
	if s.records == nil {
		s.records = make(map[int]string)
	}
	s.records[id] = text
	return nil
}

func (s *LocalRecord) ReadAt(id int) (string, error) {
	if s.records == nil {
		s.records = make(map[int]string)
	}
	if text, ok := s.records[id]; ok {
		return text, nil
	}
	return "", fmt.Errorf("not found")
}
