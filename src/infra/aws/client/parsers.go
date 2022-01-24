package client

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/consensys/quorum-key-manager/pkg/errors"
)

func parseSecretsManagerErrorResponse(err error) error {
	aerr, ok := err.(awserr.Error)
	if !ok {
		return errors.AWSError(err.Error())
	}

	switch aerr.Code() {
	case secretsmanager.ErrCodeResourceExistsException:
		return errors.AlreadyExistsError(aerr.Error())
	case secretsmanager.ErrCodeInvalidRequestException:
		return errors.InvalidFormatError(aerr.Error())
	case secretsmanager.ErrCodeResourceNotFoundException:
		return errors.NotFoundError(aerr.Error())
	case
		secretsmanager.ErrCodeInvalidNextTokenException,
		secretsmanager.ErrCodeMalformedPolicyDocumentException,
		secretsmanager.ErrCodeInvalidParameterException:
		return errors.InvalidParameterError(aerr.Error())
	default:
		return errors.AWSError(aerr.Error())
	}
}

func parseKmsErrorResponse(err error) error {
	aerr, ok := err.(awserr.Error)
	if !ok {
		return errors.AWSError(err.Error())
	}

	switch aerr.Code() {
	case kms.ErrCodeNotFoundException:
		return errors.NotFoundError(aerr.Error())
	case kms.ErrCodeAlreadyExistsException:
		return errors.AlreadyExistsError(aerr.Error())
	case
		kms.ErrCodeIncorrectKeyException,
		kms.ErrCodeIncorrectKeyMaterialException,
		kms.ErrCodeInvalidAliasNameException,
		kms.ErrCodeInvalidCiphertextException,
		kms.ErrCodeInvalidArnException:
		return errors.InvalidFormatError(aerr.Error())
	case kms.ErrCodeInvalidStateException:
		return errors.StatusConflictError(aerr.Error())
	default:
		return errors.AWSError(aerr.Error())
	}
}
