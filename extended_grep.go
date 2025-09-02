package the_platinum_searcher

import (
	"io"
	"os"
)

type extendedGrep struct {
	lineGrep
	pattern pattern
}

func (g extendedGrep) grep(path string, buf []byte) error {
	f, err := getFileHandler(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if f == os.Stdin {
		// TODO: File type is fixed in ASCII because it can not determine the character code.
		g.grepEachLines(f, ASCII, func(b []byte) bool {
			return g.pattern.regexp.Match(b)
		}, func(b []byte) int {
			return g.pattern.regexp.FindIndex(b)[0] + 1
		})
		return nil
	}

	c, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}

	if err == io.EOF {
		return nil
	}

	// detect encoding.
	limit := c
	if limit > 512 {
		limit = 512
	}

	encoding := detectEncoding(buf[:limit])
	if encoding == ERROR || encoding == BINARY {
		return nil
	}

	// grep each lines.
	g.grepEachLines(f, encoding, func(b []byte) bool {
		return g.pattern.regexp.Match(b)
	}, func(b []byte) int {
		return g.pattern.regexp.FindIndex(b)[0] + 1
	})
	return nil
}
