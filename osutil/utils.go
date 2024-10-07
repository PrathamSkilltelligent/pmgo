package osutil

import (
	"strconv"
	"syscall"

	errors "github.com/PrathamSkilltelligent/pmgo/error"
	"github.com/samber/mo"
)

func GetEnvVar(varName string) mo.Result[string] {
	val, found := syscall.Getenv(varName)
	if !found {
		return mo.Err[string](errors.EnvVarNotDefined(varName, nil))
	} else {
		return mo.Ok(val)
	}
}

func GetIntEnvVar(varName string) mo.Result[uint64] {
	v, f := GetEnvVar(varName).Get()
	if f != nil {
		return mo.Err[uint64](f)
	}

	intV, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return mo.Err[uint64](errors.EnvVarMustBeInteger(varName, err))
	}

	return mo.Ok(intV)
}
