package database

import (
	"errors"
)

const INVITE_AWAITING = 0
const INVITE_ACCEPT = 1
const INVITE_DECLINE = 2

const TIME_FORMAT = "2006-02-01 15:04:05"

var ErrINVALIDPASSWORD = errors.New("invalid password")
var ErrNoPrivilages = errors.New("insufficient privilages")
var ErrInternal = errors.New("internal transaction error")
