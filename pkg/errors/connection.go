package errors

const (
	Connection     = "CN000"
	AKV            = "CN100"
	HashicorpVault = "CN200"
	AWS            = "CN300"
	Healthcheck    = "CN400"
	BlockchainNode = "CN500"
	Postgres       = "CN600"

	InvalidRequest   = "IR000"
	Unauthorized     = "IR100"
	NotSupported     = "IR200"
	NotImplemented   = "IR300"
	InvalidFormat    = "IR400"
	InvalidParameter = "IR500"
	Forbidden        = "IR600"
)

// HashicorpVaultError is raised when failing to perform on Hashicorp Vault
func HashicorpVaultError(format string, a ...interface{}) *Error {
	return Errorf(HashicorpVault, format, a...)
}

// IsHashicorpVaultError indicate whether an error is a Hashicorp Vault connection error
func IsHashicorpVaultError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), HashicorpVault)
}

// AKVError is raised when failing to perform on AKV client
func AKVError(format string, a ...interface{}) *Error {
	return Errorf(AKV, format, a...)
}

// IsAKVError indicate whether an error is a AKV client connection error
func IsAKVError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), AKV)
}

// AWSError is raised when failing to perform on AWS client
func AWSError(format string, a ...interface{}) *Error {
	return Errorf(AWS, format, a...)
}

// IsAWSError indicate whether an error is a AWS client connection error
func IsAWSError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), AWS)
}

// PostgresError is raised when failing to perform on Postgres client
func PostgresError(format string, a ...interface{}) *Error {
	return Errorf(Postgres, format, a...)
}

// IsPostgresError indicate whether an error is a Postgres client connection error
func IsPostgresError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Postgres)
}

// HealthcheckError is raised when failing to perform a health check
func HealthcheckError(format string, a ...interface{}) *Error {
	return Errorf(Healthcheck, format, a...)
}

// IsHealthcheckError indicate whether an error is a healthcheck error
func IsHealthcheckError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Healthcheck)
}

// BlockchainNodeError is raised when failing to perform a call to the Ethereum node
func BlockchainNodeError(format string, a ...interface{}) *Error {
	return Errorf(BlockchainNode, format, a...)
}

// UnauthorizedError is raised when authentication credentials are invalid
func UnauthorizedError(format string, a ...interface{}) *Error {
	return Errorf(Unauthorized, format, a...)
}

func IsUnauthorizedError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Unauthorized)
}

// ForbiddenError is raised when the user cannot perform a given operation (missing permission)
func ForbiddenError(format string, a ...interface{}) *Error {
	return Errorf(Forbidden, format, a...)
}

func IsForbiddenError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), Forbidden)
}

// NotSupportedError is raised when operation is not supported
func NotSupportedError(format string, a ...interface{}) *Error {
	return Errorf(NotSupported, format, a...)
}

func IsNotSupportedError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), NotSupported)
}

// InvalidFormatError is raised when a Data does not match an expected format
func InvalidFormatError(format string, a ...interface{}) *Error {
	return Errorf(InvalidFormat, format, a...)
}

// IsInvalidFormatError indicate whether an error is an invalid format error
func IsInvalidFormatError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), InvalidFormat)
}

// InvalidParameterError is raised when a provided parameter invalid
func InvalidParameterError(format string, a ...interface{}) *Error {
	return Errorf(InvalidParameter, format, a...)
}

// IsInvalidParameterError indicate whether an error is an invalid parameter error
func IsInvalidParameterError(err error) bool {
	return isErrorClass(FromError(err).GetCode(), InvalidParameter)
}
