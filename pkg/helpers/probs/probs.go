package probs

var readinessErr, livenessErr error

func SetLivenessErrErr(e error) {
	livenessErr = e
}

func GetLivenessErr() error {
	return livenessErr
}

func SetReadinessErr(e error) {
	readinessErr = e
}

func GetReadinessErr() error {
	return readinessErr
}
