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
| ErrDeleteOneself | 110004 | 400 | Unable to delete yourself |
| ErrUserAlreadyDisabled | 110005 | 400 | The user is already disabled |
| ErrUserAlreadyEnabled | 110006 | 400 | The user is already enabled |
| ErrUserHasDisabled | 110007 | 401 | The account has been disabled |
| ErrReachMaxCount | 110101 | 400 | Secret reach the max count |
| ErrSecretNotFound | 110102 | 404 | Secret not found |
| ErrPolicyNotFound | 110201 | 404 | Policy not found |
| ErrPolicyAlreadyExist | 110202 | 400 | Policy already exist |
| ErrPolicyNameIsEmpty | 110203 | 400 | Policy name is empty |
| ErrResourceNotFound | 110301 | 404 | Resource not found |
| ErrResourceAlreadyExist | 110302 | 400 | Resource already exist |
| ErrResourceNameIsEmpty | 110303 | 400 | Resource name is empty |
| ErrResourceInstanceIdIsEmpty | 110304 | 400 | Resource instanceId is empty |
| ErrResourceHasAssignedPolicy | 110305 | 400 | The resource has been assigned permission policies |
| ErrRoleNotFound | 110401 | 404 | Role not found |
| ErrRoleAlreadyExist | 110402 | 400 | Role already exist |
| ErrRoleNameIsEmpty | 110403 | 400 | Role name is empty |
| ErrRoleHasAssignedUser | 110404 | 400 | The role has been assigned to a user |
| ErrAssignRoleFailed | 110405 | 400 | User role assignment fails. Please check the role status or contact the administrator |
| ErrUnsupportedAssignTarget | 110406 | 400 | The assignment target is not supported. Only user or department are supported |
| ErrRevokeRoleFailed | 110407 | 400 | User role revoke fails. Please check the role status or contact the administrator |
| ErrUnsupportedRevokeTarget | 110408 | 400 | The revoke target is not supported. Only user or department are supported |
| ErrOrgNotFound | 110501 | 404 | Organization not found |
| ErrOrgAlreadyExist | 110502 | 400 | Organization already exist |
| ErrOrgAlreadyDisabled | 110503 | 400 | The organization is already disabled |
| ErrOrgAlreadyEnabled | 110504 | 400 | The organization is already enabled |
| ErrOrgHasDisabled | 110505 | 401 | The organization has been disabled |
| ErrCannotDeleteBuiltInOrg | 110506 | 400 | Built-in organizations cannot be deleted |
| ErrCannotDisableBuiltInOrg | 110507 | 400 | Built-in organizations cannot be disabled |
| ErrMaxDepartmentsReached | 110508 | 400 | The number of departments has reached its limit |
| ErrMemberAlreadyInDepartment | 110601 | 400 | Member is already in department |
| ErrSubDepartmentsExist | 110602 | 400 | Sub departments exist and cannot be deleted |
| ErrSuccess | 100001 | 200 | OK |
| ErrUnknown | 100002 | 500 | Internal server error |
| ErrBind | 100003 | 400 | Error occurred while binding the request body to the struct |
| ErrValidation | 100004 | 400 | Validation failed |
| ErrParam | 100005 | 400 | Invalid request params |
| ErrPageNotFound | 100006 | 404 | Page not found |
| ErrInvalidRequest | 100007 | 404 | Invalid request |
| ErrDatabase | 100101 | 500 | Database error |
| ErrPrimaryEmpty | 100102 | 500 | Primary key is empty |
| ErrNilEntity | 100103 | 500 | Entity is nil |
| ErrRecordExist | 100104 | 500 | Data record is exist |
| ErrRecordNotExist | 100105 | 500 | Data record is not exist |
| ErrIndexInvalid | 100106 | 500 | Entity index is invalid |
| ErrEntityInvalid | 100107 | 500 | Entity is invalid |
| ErrTableNameEmpty | 100108 | 500 | Entity table name is empty |
| ErrEncrypt | 100201 | 401 | Error occurred while encrypting the user password |
| ErrTokenInvalid | 100202 | 401 | Token invalid |
| ErrSignatureInvalid | 100203 | 401 | Signature is invalid |
| ErrTokenMalformed | 100204 | 401 | Token is malformed |
| ErrTokenNotValidYet | 100205 | 401 | Token is not valid yet |
| ErrExpired | 100206 | 401 | Token expired |
| ErrTokenIssuedAt | 100207 | 401 | Token used before issued |
| ErrMissingLoginValues | 100208 | 401 | Missing Username or Password |
| ErrInvalidAuthHeader | 100209 | 401 | Invalid authorization header |
| ErrMissingHeader | 100210 | 401 | The `Authorization` header was empty |
| ErrPasswordIncorrect | 100211 | 401 | Invalid Username or Password |
| ErrInvalidRefreshToken | 100212 | 401 | Refresh token format is incorrect, please check |
| ErrUnauthorized | 100213 | 403 | Unauthorized |
| ErrPermissionDenied | 100214 | 403 | Permission denied |
| ErrIdentityProviderNotFound | 100215 | 401 | Identity provider not found |
| ErrEncodingFailed | 100301 | 500 | Encoding failed due to an error with the data |
| ErrDecodingFailed | 100302 | 500 | Decoding failed due to an error with the data |
| ErrInvalidJSON | 100303 | 500 | Data is not valid JSON |
| ErrEncodingJSON | 100304 | 500 | JSON data could not be encoded |
| ErrDecodingJSON | 100305 | 500 | JSON data could not be decoded |
| ErrInvalidYaml | 100306 | 500 | Data is not valid Yaml |
| ErrEncodingYaml | 100307 | 500 | Yaml data could not be encoded |
| ErrDecodingYaml | 100308 | 500 | Yaml data could not be decoded |

