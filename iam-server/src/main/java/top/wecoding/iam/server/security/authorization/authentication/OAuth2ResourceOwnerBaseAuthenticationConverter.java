package top.wecoding.iam.server.security.authorization.authentication;

import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.core.OAuth2ErrorCodes;
import org.springframework.security.oauth2.core.endpoint.OAuth2ParameterNames;
import org.springframework.security.web.authentication.AuthenticationConverter;
import org.springframework.util.MultiValueMap;
import top.wecoding.commons.core.util.JsonUtil;
import top.wecoding.commons.lang.Objects;
import top.wecoding.commons.lang.Strings;
import top.wecoding.iam.common.model.request.LoginRequest;
import top.wecoding.iam.common.util.OAuth2EndpointUtils;

import java.io.IOException;
import java.io.InputStream;
import java.util.Arrays;
import java.util.HashSet;
import java.util.Set;

/**
 * 自定义模式认证转换器
 *
 * @author liuyuhui
 * @date 2022/10/3
 */
@Slf4j
public abstract class OAuth2ResourceOwnerBaseAuthenticationConverter<
        T extends OAuth2ResourceOwnerBaseAuthenticationToken>
    implements AuthenticationConverter {

  public abstract boolean support(String authType);

  public void checkParams(LoginRequest loginRequest) {}

  public abstract T buildToken(
      Authentication clientPrincipal, Set<String> requestedScopes, LoginRequest loginRequest);

  @Override
  public Authentication convert(HttpServletRequest request) {

    // auth_type (REQUIRED)
    LoginRequest loginRequest = getLoginRequest(request);
    if (Objects.isEmpty(loginRequest) || !support(loginRequest.getAuthType())) {
      return null;
    }

    MultiValueMap<String, String> parameters = OAuth2EndpointUtils.getParameters(request);
    // scope (OPTIONAL)
    String scope = parameters.getFirst(OAuth2ParameterNames.SCOPE);
    if (Strings.hasText(scope) && parameters.get(OAuth2ParameterNames.SCOPE).size() != 1) {
      OAuth2EndpointUtils.throwError(
          OAuth2ErrorCodes.INVALID_REQUEST,
          OAuth2ParameterNames.SCOPE,
          OAuth2EndpointUtils.ACCESS_TOKEN_REQUEST_ERROR_URI);
    }

    Set<String> requestedScopes = null;
    if (Strings.hasText(scope)) {
      requestedScopes =
          new HashSet<>(Arrays.asList(Strings.delimitedListToStringArray(scope, " ")));
    }

    // verify personalization parameters
    checkParams(loginRequest);

    // obtain currently authenticated client
    Authentication clientPrincipal = SecurityContextHolder.getContext().getAuthentication();
    if (clientPrincipal == null) {
      OAuth2EndpointUtils.throwError(
          OAuth2ErrorCodes.INVALID_REQUEST,
          OAuth2ErrorCodes.INVALID_CLIENT,
          OAuth2EndpointUtils.ACCESS_TOKEN_REQUEST_ERROR_URI);
    }

    // build token
    return buildToken(clientPrincipal, requestedScopes, loginRequest);
  }

  private LoginRequest getLoginRequest(HttpServletRequest request) {
    try (InputStream inputStream = request.getInputStream()) {
      return JsonUtil.readValue(inputStream, LoginRequest.class);
    } catch (IOException e) {
      log.warn("Failed to obtain login parameters {}.", e.getMessage());
    }
    return null;
  }
}
