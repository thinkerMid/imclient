package external

import "errors"

// TODO 归类底层产生的错误 提供上层判断

// AccountNotFoundErr .
var AccountNotFoundErr = errors.New("account not found")
