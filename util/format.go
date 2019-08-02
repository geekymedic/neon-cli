package util

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/geekymedic/neon/logger"
	"github.com/kyokomi/emoji"
)

func StdoutExit(exitCode int, format string, args ...interface{}) {
	if exitCode == 0 {
		emoji.Fprintf(os.Stdout, ":heavy_check_mark:"+format, args...)
	} else {
		emoji.Fprintf(os.Stdout, ":heavy_multiplication_x:"+format, args...)
	}
	os.Exit(exitCode)
}

func StdDebug(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func StdoutOk(format string, args ...interface{}) {
	emoji.Fprintf(os.Stdout, ":heavy_check_mark:"+format, args...)
}

func neonPic() string {
	txt := `

	\\		  //
	 \\		 //
	  \\	//	
	  //--- \\	
	 //		 \\
	//		  \\

http://www.geekymedic.cn
`
	return txt
}

func ParseSystemName(dir string) (string, string) {
	pathSlice := strings.Split(dir, string(filepath.Separator))
	systemName := pathSlice[len(pathSlice)-1]
	idx := strings.Index(systemName, "_")
	if idx < 0 {
		StdoutExit(-1, "%s is invalid service name", systemName)
	}
	systemName = systemName[:idx]
	return systemName, pathSlice[len(pathSlice)-1]
}

func ProcessBar(name string, position, count int, upRate ...float32) float32 {
	if position > count {
		return 1
	}
	total := float32(80)
	rate := float32(position) / float32(count)
	leftProcess := int(total * rate)
	rightProcess := 80 - leftProcess
	rate = float32(int(rate*100)) / 100
	if len(upRate) > 0 && upRate[0] >= rate {
		return rate
	}
	StdoutOk(":earth_africa:%s:%s%s\n", name, strings.Repeat("=", leftProcess), strings.Repeat("-", rightProcess))
	return rate
}
