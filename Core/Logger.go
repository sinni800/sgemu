package Core

import (
	"fmt"
	"io"
	"log"
	"runtime"
)

const (
	LOG_DEBUG   = 1
	LOG_WARNING = 2
	LOG_INFO    = 4
	LOG_ALL     = LOG_DEBUG | LOG_WARNING | LOG_INFO
)

type multiWriter struct {
	writers []io.Writer
}

func (t *multiWriter) Write(p []byte) (n int, err error) {
	for i := 0; i < len(t.writers); i++ {
		n, e := t.writers[i].Write(p)
		if e != nil {
			t.writers[len(t.writers)-1], t.writers[i], t.writers = nil, t.writers[len(t.writers)-1], t.writers[:len(t.writers)-1]
			err = e
			i--
			continue
		}
		if n != len(p) {
			t.writers[len(t.writers)-1], t.writers[i], t.writers = nil, t.writers[len(t.writers)-1], t.writers[:len(t.writers)-1]
			err = io.ErrShortWrite
			i--
			continue
		}
	}
	return len(p), err
}

func (t *multiWriter) AddWriter(out io.Writer) {
	t.writers = append(t.writers, out)
}

func MultiWriter(writers ...io.Writer) *multiWriter {
	w := &multiWriter{make([]io.Writer, 0, 10)}
	for _, writer := range writers {
		w.writers = append(w.writers, writer)
	}
	return w
}

type Logger struct {
	log.Logger
	out   *multiWriter
	level byte
}

func (l *Logger) AddWriter(out io.Writer) {
	l.out.AddWriter(out)
}

func NewLogger(out io.Writer, prefix string, flag int) *Logger {
	l := &Logger{}
	l.out = MultiWriter(out)
	l.Logger = *log.New(l.out, prefix, flag)
	l.level = LOG_ALL
	return l
}

func (l *Logger) SetLogLevel(level byte) {
	l.level = level
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	panic(s)
}

func (l *Logger) Println_Debug(v ...interface{}) {
	if l.level&LOG_DEBUG > 0 {
		l.Println(v...)
	}
}

func (l *Logger) Printf_Debug(format string, v ...interface{}) {
	if l.level&LOG_DEBUG > 0 {
		l.Printf(format, v...)
	}
}

func (l *Logger) Print_Debug(v ...interface{}) {
	if l.level&LOG_DEBUG > 0 {
		l.Print(v...)
	}
}

func (l *Logger) Println_Warning(v ...interface{}) {
	if l.level&LOG_WARNING > 0 {
		l.Println(v...)
	}
}

func (l *Logger) Printf_Warning(format string, v ...interface{}) {
	if l.level&LOG_WARNING > 0 {
		l.Printf(format, v...)
	}
}

func (l *Logger) Print_Warning(v ...interface{}) {
	if l.level&LOG_WARNING > 0 {
		l.Print(v...)
	}
}

func (l *Logger) Println_Info(v ...interface{}) {
	if l.level&LOG_INFO > 0 {
		l.Println(v...)
	}
}

func (l *Logger) Printf_Info(format string, v ...interface{}) {
	if l.level&LOG_INFO > 0 {
		l.Printf(format, v...)
	}
}

func (l *Logger) Print_Info(v ...interface{}) {
	if l.level&LOG_INFO > 0 {
		l.Print(v...)
	}
}

func PanicPath() string {
	fullPath := ""
	skip := 3
	for i := skip; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if i > skip {
			fullPath += ", "
		}
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		fullPath += fmt.Sprintf("%s:%d", file, line)
	}
	return fullPath
}
