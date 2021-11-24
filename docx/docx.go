package docx

import (
	"github.com/nguyenthenguyen/docx"
	log "github.com/sirupsen/logrus"
)

func EditDoc(filein, fileout string, oldText, newText []string) bool {
	log.Info(filein)
	log.Info(fileout)
	r, err := docx.ReadDocxFile(filein)
	// Or read from memory
	// r, err := docx.ReadDocxFromMemory(data io.ReaderAt, size int64)
	if err != nil {
		log.Error(err)
		return false
	}

	docx1 := r.Editable()

	for i := range oldText {
		log.Info(oldText[i], " || ", newText[i])
		docx1.Replace(oldText[i], newText[i], -1)
	}

	docx1.WriteToFile(fileout)
	r.Close()

	return true
}
