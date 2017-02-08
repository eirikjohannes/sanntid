package hardware

/*
#cgo LDFLAGS: -lcomedi -lm
#include "channels.h"
*/
import "C"
import (
        def "definitions"
        "log"
	"elev"
)

elev.Init()
