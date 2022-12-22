package panicDump

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"runtime"
	"time"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// PanicProfile .
type PanicProfile struct {
	DumpID     string
	Time       time.Time
	StackFrame []string
}

// Print .
func (p PanicProfile) Print(log *zap.SugaredLogger) {
	log = log.Named("PanicDump")

	for i := range p.StackFrame {
		log.Error(p.StackFrame[i])
	}
}

func newPanicProfile() PanicProfile {
	id := make([]byte, 12)
	_, _ = rand.Read(id)

	dumpID := hex.EncodeToString(id)
	start := time.Now()

	return PanicProfile{
		DumpID: dumpID,
		Time:   start,
	}
}

// Scan formatted stack frame, skipping skip frames. returns a dumpID
func Scan(skip int) PanicProfile {
	p := newPanicProfile()

	p.StackFrame = append(p.StackFrame, fmt.Sprintf("%s  runtime caller stack:", p.DumpID))

	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	var lineIndex uint8
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		lineIndex++
		p.StackFrame = append(p.StackFrame, fmt.Sprintf("%s\t%v %s:%d (0x%x)", p.DumpID, lineIndex, file, line, pc))

		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}

		lineIndex++
		p.StackFrame = append(p.StackFrame, fmt.Sprintf("%s\t%v \t%s: %s", p.DumpID, lineIndex, function(pc), source(lines, line)))
	}

	return p
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in Scan trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
