package client

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
)

func parseSecretsManagerErrorResponse(err error) error {
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

func parseKmsErrorResponse(err error) error {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case kms.ErrCodeAlreadyExistsException:
			return errors.AlreadyExistsError("resource already exists")
		case kms.ErrCodeInternalException:
			return errors.AWSError("internal error")
		case kms.ErrCodeLimitExceededException:
			return errors.AWSError("resource limit error")
		case kms.ErrCodeIncorrectKeyException:
		case kms.ErrCodeIncorrectKeyMaterialException:
		case kms.ErrCodeInvalidAliasNameException:
		case kms.ErrCodeInvalidCiphertextException:
		case kms.ErrCodeInvalidArnException:
		case kms.ErrCodeInvalidStateException:
			return errors.InvalidParameterError("invalid parameter")
		case kms.ErrCodeNotFoundException:
			return errors.NotFoundError("resource was not found")

		}
	}
	return err
}
