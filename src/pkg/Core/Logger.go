package Core

import( 
	"runtime"
	 "fmt"
	 "log"
	 "io"
	 "container/list"
 )
 
type multiWriter struct {
	writers *list.List
}

func (t *multiWriter) Write(p []byte) (n int, err error) {
	for e := t.writers.Front(); e != nil; e = e.Next() {
		w := e.Value.(io.Writer)
		n, err = w.Write(p)
		if err != nil {
			t.writers.Remove(e)
			continue
		}
		if n != len(p) {
			err = io.ErrShortWrite
			t.writers.Remove(e)
			continue
		}
	}
	return len(p), nil
}

func (t *multiWriter) AddWriter(out io.Writer) {
	t.writers.PushBack(out)
}

func MultiWriter(writers ...io.Writer) *multiWriter {
	l := list.New()
	for _,writer := range writers {
		l.PushBack(writer)
	}
	return &multiWriter{l}
}

type Logger struct {
	log.Logger
	out *multiWriter
}

func (l *Logger) AddWriter(out io.Writer) {
 	l.out.AddWriter(out)
}

func NewLogger(out io.Writer, prefix string, flag int) *Logger {
	l := &Logger{}
	l.out = MultiWriter(out)
	l.Logger = *log.New(l.out, prefix, flag)
	return l
}

func PanicPath() string{
	fullPath := ""
	skip := 3
	for i:=skip;;i++ {
		_, file, line, ok := runtime.Caller(i)
   		if !ok {
   				break;
   		}
   		if i > skip {
   			fullPath += ", ";
   		}
		short := file
  		for i := len(file) - 1; i > 0; i-- {
  				if file[i] == '/' {
   					short = file[i+1:]
   					break
  				} 
   		}
	   	file = short
	   	fullPath += fmt.Sprintf("%s:%d", file , line)
   	}
   	return fullPath
}  

func (l *Logger) Panicf(format string, v ...interface{}) {
		s := fmt.Sprintf(format, v...)
		panic(s)
}