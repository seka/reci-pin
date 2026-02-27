package errors

import "errors"

var (
	// ErrNotFound は、要求されたリソースが見つからない場合に返されます。
	ErrNotFound = errors.New("resource not found")

	// ErrAlreadyExists は、リソースが既に存在する場合に返されます。
	ErrAlreadyExists = errors.New("resource already exists")

	// ErrBadRequest は、リクエストデータが不正な場合に返されます。
	ErrBadRequest = errors.New("bad request")

	// ErrIntegrityConstraint は、データベースの整合性制約に違反した場合に返されます。
	ErrIntegrityConstraint = errors.New("integrity constraint violation")

	// ErrTransaction は、トランザクション関連のエラーが発生した場合に返されます。
	ErrTransaction = errors.New("transaction error")

	// ErrInternal は、予期しない内部エラーが発生した場合に返されます。
	ErrInternal = errors.New("internal server error")
)
