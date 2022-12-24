package top.wecoding.iam.server.security.handler;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.time.temporal.ChronoUnit;
import java.util.List;
import java.util.Map;
import lombok.SneakyThrows;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.converter.HttpMessageConverter;
import org.springframework.http.server.ServletServerHttpResponse;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.security.oauth2.core.OAuth2RefreshToken;
import org.springframework.security.oauth2.core.endpoint.OAuth2AccessTokenResponse;
import org.springframework.security.oauth2.core.http.converter.OAuth2AccessTokenResponseHttpMessageConverter;
import org.springframework.security.oauth2.server.authorization.authentication.OAuth2AccessTokenAuthenticationToken;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenClaimNames;
import org.springframework.security.web.authentication.AuthenticationSuccessHandler;
import org.springframework.util.CollectionUtils;
import top.wecoding.commons.lang.Objects;
import top.wecoding.iam.server.util.LogUtil;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
@Slf4j
@SuppressWarnings("unchecked")
public class SsoAuthenticationSuccessHandler implements AuthenticationSuccessHandler {

  private final HttpMessageConverter<OAuth2AccessTokenResponse> accessTokenHttpResponseConverter =
      new OAuth2AccessTokenResponseHttpMessageConverter();

  /**
   * Called when a user has been successfully authenticated.
   *
   * @param request the request which caused the successful authentication
   * @param response the response
   * @param authentication the <tt>Authentication</tt> object which was created during the
   *     authentication process.
   */
  @Override
  @SneakyThrows
  public void onAuthenticationSuccess(
      HttpServletRequest request, HttpServletResponse response, Authentication authentication) {
    OAuth2AccessTokenAuthenticationToken accessTokenAuthentication =
        (OAuth2AccessTokenAuthenticationToken) authentication;
    Map<String, Object> map = accessTokenAuthentication.getAdditionalParameters();
    if (!Objects.isEmpty(map) && map.containsKey(OAuth2TokenClaimNames.SUB)) {
      List<String> audience = (List<String>) map.get(OAuth2TokenClaimNames.AUD);
      String userId = (String) map.get(OAuth2TokenClaimNames.SUB);
      log.info("user {} login successful.", audience);
      SecurityContextHolder.getContext().setAuthentication(accessTokenAuthentication);
      LogUtil.successLogin(userId);
    }

    sendAccessTokenResponse(request, response, authentication);
  }

  private void sendAccessTokenResponse(
      HttpServletRequest request, HttpServletResponse response, Authentication authentication)
      throws IOException {

    OAuth2AccessTokenAuthenticationToken accessTokenAuthentication =
        (OAuth2AccessTokenAuthenticationToken) authentication;

    OAuth2AccessToken accessToken = accessTokenAuthentication.getAccessToken();
    OAuth2RefreshToken refreshToken = accessTokenAuthentication.getRefreshToken();
    Map<String, Object> additionalParameters = accessTokenAuthentication.getAdditionalParameters();

    OAuth2AccessTokenResponse.Builder builder =
        OAuth2AccessTokenResponse.withToken(accessToken.getTokenValue())
            .tokenType(accessToken.getTokenType())
            .scopes(accessToken.getScopes());
    if (accessToken.getIssuedAt() != null && accessToken.getExpiresAt() != null) {
      builder.expiresIn(
          ChronoUnit.SECONDS.between(accessToken.getIssuedAt(), accessToken.getExpiresAt()));
    }
    if (refreshToken != null) {
      builder.refreshToken(refreshToken.getTokenValue());
    }
    if (!CollectionUtils.isEmpty(additionalParameters)) {
      builder.additionalParameters(additionalParameters);
    }
    OAuth2AccessTokenResponse accessTokenResponse = builder.build();
    ServletServerHttpResponse httpResponse = new ServletServerHttpResponse(response);

    // 无状态 注意删除 context 上下文的信息
    SecurityContextHolder.clearContext();
    this.accessTokenHttpResponseConverter.write(accessTokenResponse, null, httpResponse);
  }
}
