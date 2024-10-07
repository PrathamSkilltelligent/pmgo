package errors

import "github.com/PrathamSkilltelligent/pmgo/fault"

const (
	component fault.ErrComponent = "lyra.go"
)
const (
	ErrEnvVarUndefined    fault.ErrorCode = "LYRAGO000000"
	ErrInvalidEnvVarType  fault.ErrorCode = "LYRAGO000005"
	ErrCreatingSmtpClient fault.ErrorCode = "LYRAGO000010"
)

func EnvVarNotDefined(varName string, cause error) fault.Fault {
	data := map[string]any{"var_name": varName}
	return fault.NewBasicFault(ErrEnvVarUndefined).
		SetComponent(component).
		ToFault(data, cause)
}
func EnvVarMustBeInteger(varName string, cause error) fault.Fault {
	data := map[string]any{"var_name": varName, "var_type": "integer"}
	return fault.NewBasicFault(ErrInvalidEnvVarType).
		SetComponent(component).
		ToFault(data, cause)
}

func SmtpClientCreationError(cause error) fault.Fault {
	return fault.NewBasicFault(ErrCreatingSmtpClient).
		SetComponent(component).
		ToFault(nil, cause)
}
