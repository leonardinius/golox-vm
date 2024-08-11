package vmdebug

import "os"

type PprofDebug struct {
	On         bool
	CPUProfile string
	MemProfile string
}

func LoadPprofConfigFromEnv() *PprofDebug {
	isOn := func(s string) bool {
		return s != "" && s != "0" && s != "false" && s != "no"
	}
	isOff := func(s string) bool {
		return s == "0" || s == "false" || s == "no"
	}
	valueOrDefault := func(value, def string) string {
		if value == "" {
			return def
		}
		return value
	}

	profile := &PprofDebug{}
	if on := os.Getenv("GLOX_PPROF"); isOn(on) {
		profile.On = true
		if cpuOn := os.Getenv("GLOX_PPROF_CPU"); !isOff(cpuOn) {
			profile.CPUProfile = valueOrDefault(os.Getenv("GLOX_PPROF_CPU_NAME"), "cpu.pb.gz")
		}
		if memOn := os.Getenv("GLOX_PPROF_MEM"); !isOff(memOn) {
			profile.MemProfile = valueOrDefault(os.Getenv("GLOX_PPROF_MEM_NAME"), "mem.pb.gz")
		}
	}

	return profile
}
