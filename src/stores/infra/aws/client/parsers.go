package client

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func parseErrorResponse(err error) error {
	aerr, ok := err.(awserr.Error)
	if !ok {
		return errors.AWSError(err.Error())
	}

	switch aerr.Code() {
	case secretsmanager.ErrCodeResourceExistsException:
		return errors.AlreadyExistsError(aerr.Error())
	case secretsmanager.ErrCodeInvalidParameterException:
		return errors.InvalidParameterError(aerr.Error())
	case secretsmanager.ErrCodeInvalidRequestException:
		return errors.InvalidFormatError(aerr.Error())
	case secretsmanager.ErrCodeResourceNotFoundException:
		return errors.NotFoundError(aerr.Error())
	case secretsmanager.ErrCodeInvalidNextTokenException:
		return errors.InvalidParameterError(aerr.Error())
	case secretsmanager.ErrCodeMalformedPolicyDocumentException:
		return errors.InvalidParameterError(aerr.Error())
	default:
		return errors.AWSError(aerr.Error())
	}
}
