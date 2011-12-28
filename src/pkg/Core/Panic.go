package Core

import "runtime"
import "fmt"

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

func Spanicf(format string, v ...interface{}) {
		s := fmt.Sprintf(format, v...)
		panic(s)
}