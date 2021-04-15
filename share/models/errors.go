package models

import "errors"

var (
	ErrMessageDeleteProhibitedWithAffix = errors.New("message: Prohibit delete messages with attachments")
	ErrNotAllowedToDelete               = errors.New("Not allowed to delete")
	ErrNotAllowedToUpdate               = errors.New("Not allowed to update")
	ErrNotFoundMemberID                 = errors.New("Not found MemeberID")
	ErrParameterError                   = errors.New("Parameter error")
	ErrParameterErrorWithExchange       = errors.New("exchange: Parameter error")
	ErrUndefinedMemberID                = errors.New("Undefined MemberID property")
	ErrUndefinedTeamID                  = errors.New("Undefined TeamID property")
	ErrUploadCountOutRange              = errors.New("upload: File Count is out of range")
	ErrUploadSizeOutRange               = errors.New("upload: File size is out of range")
	ErrValidateCreator                  = errors.New("validate: error ")
	ErrMembeber                         = errors.New("Member data error")
	ErrNotFoundMessageID                = errors.New("Message : not found message")
)
