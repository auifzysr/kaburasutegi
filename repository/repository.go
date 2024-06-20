package repository

type Recorder interface {
	Record(text string) error
	RecordAt(id int, text string) error
	ReadAt(id int) (string, error)
}
