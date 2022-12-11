package top.wecoding.iam.server.security.authorization.authentication.password;

import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.core.OAuth2ErrorCodes;
import org.springframework.security.oauth2.core.endpoint.OAuth2ParameterNames;
import top.wecoding.commons.lang.Strings;
import top.wecoding.iam.common.constant.AuthParameterNames;
import top.wecoding.iam.common.enums.AuthType;
import top.wecoding.iam.common.model.request.LoginRequest;
import top.wecoding.iam.common.util.OAuth2EndpointUtils;
import top.wecoding.iam.server.security.authorization.authentication.OAuth2ResourceOwnerBaseAuthenticationConverter;

import java.util.Objects;
import java.util.Set;

/**
 * 密码认证转换器
 *
 * @author liuyuhui
 * @date 2022/10/3
 */
public class OAuth2ResourceOwnerPasswordAuthenticationConverter
    extends OAuth2ResourceOwnerBaseAuthenticationConverter<
            OAuth2ResourceOwnerPasswordAuthenticationToken> {

  @Override
  public boolean support(String authType) {
    return AuthType.PASSWORD.code().equalsIgnoreCase(authType);
  }

  @Override
  @SuppressWarnings({"unchecked", "rawtypes"})
  public OAuth2ResourceOwnerPasswordAuthenticationToken buildToken(
      Authentication clientPrincipal, Set requestedScopes, LoginRequest loginRequest) {
    return new OAuth2ResourceOwnerPasswordAuthenticationToken(
        AuthType.PASSWORD, clientPrincipal, requestedScopes, loginRequest);
  }

  @Override
  public void checkParams(LoginRequest loginRequest) {
    // password_payload (REQUIRED)
    if (Objects.isNull(loginRequest.getPasswordPayload())) {
      OAuth2EndpointUtils.throwError(
          OAuth2ErrorCodes.INVALID_REQUEST,
          AuthParameterNames.PASSWORD_PAYLOAD,
          OAuth2EndpointUtils.ACCESS_TOKEN_REQUEST_ERROR_URI);
    }

    // username (REQUIRED)
    String username = loginRequest.getPasswordPayload().getAccount();
    if (!Strings.hasText(username)) {
      OAuth2EndpointUtils.throwError(
          OAuth2ErrorCodes.INVALID_REQUEST,
          OAuth2ParameterNames.USERNAME,
          OAuth2EndpointUtils.ACCESS_TOKEN_REQUEST_ERROR_URI);
    }

    // password (REQUIRED)
    String password = loginRequest.getPasswordPayload().getPassword();
    if (!Strings.hasText(password)) {
      OAuth2EndpointUtils.throwError(
          OAuth2ErrorCodes.INVALID_REQUEST,
          OAuth2ParameterNames.PASSWORD,
          OAuth2EndpointUtils.ACCESS_TOKEN_REQUEST_ERROR_URI);
    }
  }
}
