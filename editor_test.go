package editor_test

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
)

// TextEditor provides functionality for manipulating and analyzing text documents.
type TextEditor interface {
	// Removes characters [i..j) from the document and places them in the clipboard.
	// Previous clipboard contents is overwritten.
	Cut(i, j int)
	// Places characters [i..j) from the document in the clipboard.
	// Previous clipboard contents is overwritten.
	Copy(i, j int)
	// Inserts the contents of the clipboard into the document starting at position i.
	// Nothing is inserted if the clipboard is empty.
	Paste(i int)
	// Returns the document as a string.
	GetText() string
	// Returns the number of misspelled words in the document. A word is considered misspelled
	// if it does not appear in /usr/share/dict/words or any other dictionary (of comparable size)
	// that you choose.
	Misspellings() int
}

type SimpleEditor struct {
	document   string
	dictionary map[string]bool
	pasteText  string
}

func NewSimpleEditor(document string) TextEditor {
	// On windows, the dictionary can often be found at:
        // C:/Users/{username}/AppData/Roaming/Microsoft/Spelling/en-US/default.dic
	fileHandle, _ := os.Open("/usr/share/dict/words")
	defer fileHandle.Close()
	dict := make(map[string]bool)
	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		dict[scanner.Text()] = true
	}
	return &SimpleEditor{document: document, dictionary: dict}
}

func (s *SimpleEditor) Cut(i, j int) {
	s.pasteText = s.document[i:j]
	s.document = s.document[:i] + s.document[j:]
}

func (s *SimpleEditor) Copy(i, j int) {
	s.pasteText = s.document[i:j]
}

func (s *SimpleEditor) Paste(i int) {
	s.document = s.document[:i] + s.pasteText + s.document[i:]
}

func (s *SimpleEditor) GetText() string {
	return s.document
}

func (s *SimpleEditor) Misspellings() int {
	result := 0
	for _, word := range strings.Fields(s.document) {
		if !s.dictionary[word] {
			result++
		}
	}
	return result
}

func BenchmarkClipboard(b *testing.B) {
	cases := []struct {
		data string
	}{
		{strings.Repeat("Neeva is awesome!", 10)},
		{strings.Repeat("Neeva is awesome!", 100)},
	}
	for _, tc := range cases {
		s := NewSimpleEditor(tc.data)
		b.Run("CutPaste"+strconv.Itoa(len(tc.data)), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				if n%2 == 0 {
					s.Cut(1, 3)
				} else {
					s.Paste(2)
				}
			}
		})
		s = NewSimpleEditor(tc.data)
		b.Run("CopyPaste"+strconv.Itoa(len(tc.data)), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				if n%2 == 0 {
					s.Copy(1, 3)
				} else {
					s.Paste(2)
				}
			}
		})
		s = NewSimpleEditor(tc.data)
		b.Run("GetText"+strconv.Itoa(len(tc.data)), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = s.GetText()
			}
		})
		s = NewSimpleEditor(tc.data)
		b.Run("Misspellings"+strconv.Itoa(len(tc.data)), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_ = s.Misspellings()
			}
		})
	}
}
