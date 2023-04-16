# 错误码

！！IAM 系统错误码列表，由 `codegen -type=int -doc` 命令生成，不要对此文件做任何更改。

## 功能说明

如果返回结果中存在 `code` 字段，则表示调用 API 接口失败。例如：

```json
{
  "code": 100101,
  "msg": "Database error"
}
```

上述返回中 `code` 表示错误码，`message` 表示该错误的具体信息。每个错误同时也对应一个 HTTP 状态码，比如上述错误码对应了 HTTP 状态码 500(Internal Server Error)。

## 错误码列表

IAM 系统支持的错误码列表如下：

| Identifier | Code | HTTP Code | Description |
| ---------- | ---- | --------- | ----------- |
| ErrUserNotFound | 110001 | 404 | User not found |
| ErrUserAlreadyExist | 110002 | 400 | User already exist |
| ErrUserNameIsEmpty | 110003 | 400 | Username is empty |
| ErrReachMaxCount | 110101 | 400 | Secret reach the max count |
| ErrSecretNotFound | 110102 | 404 | Secret not found |
| ErrPolicyNotFound | 110201 | 404 | Policy not found |
| ErrSuccess | 100001 | 200 | OK |
| ErrUnknown | 100002 | 500 | Internal server error |
| ErrBind | 100003 | 400 | Error occurred while binding the request body to the struct |
| ErrValidation | 100004 | 400 | Validation failed |
| ErrTokenInvalid | 100005 | 401 | Token invalid |
| ErrPageNotFound | 100006 | 404 | Page not found |
| ErrDatabase | 100101 | 500 | Database error |
| ErrPrimaryEmpty | 100102 | 500 | Primary key is empty |
| ErrNilEntity | 100103 | 500 | Entity is nil |
| ErrRecordExist | 100104 | 500 | Entity primary key is exist |
| ErrRecordNotExist | 100105 | 500 | Entity primary key is not exist |
| ErrIndexInvalid | 100106 | 500 | Entity index is invalid |
| ErrEntityInvalid | 100107 | 500 | Entity is invalid |
| ErrEncrypt | 100201 | 401 | Error occurred while encrypting the user password |
| ErrSignatureInvalid | 100202 | 401 | Signature is invalid |
| ErrTokenMalformed | 100203 | 401 | Token is malformed |
| ErrTokenNotValidYet | 100204 | 401 | Token is not valid yet |
| ErrExpired | 100205 | 401 | Token expired |
| ErrMissingLoginValues | 100206 | 401 | Missing Username or Password |
| ErrInvalidAuthHeader | 100207 | 401 | Invalid authorization header |
| ErrMissingHeader | 100208 | 401 | The `Authorization` header was empty |
| ErrPasswordIncorrect | 100209 | 401 | Invalid Username or Password |
| ErrPermissionDenied | 100210 | 403 | Permission denied |
| ErrEncodingFailed | 100301 | 500 | Encoding failed due to an error with the data |
| ErrDecodingFailed | 100302 | 500 | Decoding failed due to an error with the data |
| ErrInvalidJSON | 100303 | 500 | Data is not valid JSON |
| ErrEncodingJSON | 100304 | 500 | JSON data could not be encoded |
| ErrDecodingJSON | 100305 | 500 | JSON data could not be decoded |
| ErrInvalidYaml | 100306 | 500 | Data is not valid Yaml |
| ErrEncodingYaml | 100307 | 500 | Yaml data could not be encoded |
| ErrDecodingYaml | 100308 | 500 | Yaml data could not be decoded |

