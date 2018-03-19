package sandpile

type Config struct {
	NumGrains  int
	MaxGrains  int
	Colors     ColorList
	NumWorkers int

	FileName string
	FilePath string
}
