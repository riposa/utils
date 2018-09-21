package log

import (
	"fmt"
	"github.com/riposa/utils/errors"
	"github.com/fatih/color"
	"io"
	"os"
	"runtime"
	"time"
)

const (
	Ldate = 1 << iota
	Ltime
	Lmicroseconds
	Llongfile
	Lshortfile
	LUTC
	LstdFlags = Ldate | Ltime

	MaxOutput = 3

	debug     = 10
	info      = 20
	warning   = 30
	errorLv   = 40
	exception = 50
)

type logger struct {
	level int
	flag  int
	out   []io.Writer
	buf   []byte
}

var (
	infoPrefix      = color.HiGreenString("[INFO]   ")
	debugPrefix     = color.HiCyanString("[DEBUG]  ")
	warningPrefix   = color.HiYellowString("[WARNING]")
	errorPrefix     = color.HiRedString("[ERROR]  ")
	exceptionPrefix = color.HiRedString("[EXCEPTION]")
)

func New(out ...io.Writer) *logger {
	var writers []io.Writer
	if len(out) > 0 {
		writers = out
	} else {
		writers = []io.Writer{os.Stdout}
	}
	return &logger{
		level: info,
		flag:  Ldate | Ltime | Lmicroseconds | Lshortfile,
		out:   writers,
	}
}

func (l *logger) SetOutput(w io.Writer) bool {
	if len(l.out) >= MaxOutput {
		return false
	}
	l.out = append(l.out, w)
	return true
}

func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l *logger) formatHeader(buf *[]byte, t time.Time, file string, line int, prefix string) {
	*buf = append(*buf, prefix...)
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
}

func (l *logger) Output(calldepth int, s string, lv int, prefix string) error {
	now := time.Now()
	var file string
	var line int

	if lv < l.level {
		return nil
	}

	if l.flag&(Lshortfile|Llongfile) != 0 {
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line, prefix)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	for _, out := range l.out {
		_, err := out.Write(l.buf)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

func (l *logger) EnableDebug() {
	l.level = debug
}

func (l *logger) DisableDebug() {
	l.level = info
}

func (l *logger) Print(v ...interface{}) {
	l.Output(2, fmt.Sprint(v...), info, "")
}

func (l *logger) Info(v ...interface{}) {
	l.Output(2, fmt.Sprint(v...), info, infoPrefix)
}

func (l *logger) Warning(v ...interface{}) {
	l.Output(2, fmt.Sprint(v...), warning, warningPrefix)
}

func (l *logger) Debug(v ...interface{}) {
	l.Output(2, fmt.Sprint(v...), debug, debugPrefix)
}

func (l *logger) Error(v ...interface{}) {
	l.Output(2, fmt.Sprint(v...), errorLv, errorPrefix)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...), info, infoPrefix)
}

func (l *logger) Warningf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...), warning, warningPrefix)
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...), debug, debugPrefix)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...), errorLv, errorPrefix)
}

func (l *logger) Exception(err error) {
	if e, ok := err.(errors.Error); ok {
		l.Output(2, "\n"+fmt.Sprint(color.HiRedString(e.Error()))+"\n", exception, exceptionPrefix)
	} else {
		l.Output(2, fmt.Sprint(err), exception, exceptionPrefix)
	}
}
