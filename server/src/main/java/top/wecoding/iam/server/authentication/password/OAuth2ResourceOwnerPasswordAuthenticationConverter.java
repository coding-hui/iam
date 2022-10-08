package top.wecoding.iam.server.authentication.password;

import java.util.Map;
import java.util.Set;
import javax.servlet.http.HttpServletRequest;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.core.AuthorizationGrantType;
import org.springframework.security.oauth2.core.OAuth2ErrorCodes;
import org.springframework.security.oauth2.core.endpoint.OAuth2ParameterNames;
import org.springframework.util.MultiValueMap;
import org.springframework.util.StringUtils;
import top.wecoding.iam.common.util.OAuth2EndpointUtils;
import top.wecoding.iam.server.authentication.base.OAuth2ResourceOwnerBaseAuthenticationConverter;

/**
 * 密码认证转换器
 *
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
public class OAuth2ResourceOwnerPasswordAuthenticationConverter
    extends OAuth2ResourceOwnerBaseAuthenticationConverter<
        OAuth2ResourceOwnerPasswordAuthenticationToken> {

  /**
   * 支持密码模式
   *
   * @param grantType 授权类型
   */
  @Override
  public boolean support(String grantType) {
    return AuthorizationGrantType.PASSWORD.getValue().equals(grantType);
  }

  @Override
  public OAuth2ResourceOwnerPasswordAuthenticationToken buildToken(
      Authentication clientPrincipal, Set requestedScopes, Map additionalParameters) {
    return new OAuth2ResourceOwnerPasswordAuthenticationToken(
        AuthorizationGrantType.PASSWORD, clientPrincipal, requestedScopes, additionalParameters);
  }

  /**
   * 校验扩展参数 密码模式密码必须不为空
   *
   * @param request 参数列表
   */
  @Override
  public void checkParams(HttpServletRequest request) {
    MultiValueMap<String, String> parameters = OAuth2EndpointUtils.getParameters(request);
    // username (REQUIRED)
    String username = parameters.getFirst(OAuth2ParameterNames.USERNAME);
    if (!StringUtils.hasText(username)
        || parameters.get(OAuth2ParameterNames.USERNAME).size() != 1) {
      OAuth2EndpointUtils.throwError(
          OAuth2ErrorCodes.INVALID_REQUEST,
          OAuth2ParameterNames.USERNAME,
          OAuth2EndpointUtils.ACCESS_TOKEN_REQUEST_ERROR_URI);
    }

    // password (REQUIRED)
    String password = parameters.getFirst(OAuth2ParameterNames.PASSWORD);
    if (!StringUtils.hasText(password)
        || parameters.get(OAuth2ParameterNames.PASSWORD).size() != 1) {
      OAuth2EndpointUtils.throwError(
          OAuth2ErrorCodes.INVALID_REQUEST,
          OAuth2ParameterNames.PASSWORD,
          OAuth2EndpointUtils.ACCESS_TOKEN_REQUEST_ERROR_URI);
    }
  }
}
