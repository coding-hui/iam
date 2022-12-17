package top.wecoding.iam.common.convert;

import org.springframework.core.convert.converter.Converter;
import org.springframework.security.oauth2.core.endpoint.OAuth2AccessTokenResponse;
import org.springframework.security.oauth2.core.oidc.endpoint.OidcParameterNames;
import top.wecoding.commons.core.model.R;
import top.wecoding.commons.core.util.JsonUtil;
import top.wecoding.iam.common.model.response.LoginResponse;

import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.Map;

/**
 * @author liuyuhui
 * @since 0.5
 */
public class RestAccessTokenResponseHttpMessageConverter
    implements Converter<OAuth2AccessTokenResponse, Map<String, Object>> {

  private static long getExpiresIn(OAuth2AccessTokenResponse tokenResponse) {
    if (tokenResponse.getAccessToken().getExpiresAt() != null) {
      return ChronoUnit.SECONDS.between(
          Instant.now(), tokenResponse.getAccessToken().getExpiresAt());
    }
    return -1;
  }

  @Override
  @SuppressWarnings("unchecked")
  public Map<String, Object> convert(OAuth2AccessTokenResponse tokenResponse) {
    String idToken =
        (String) tokenResponse.getAdditionalParameters().get(OidcParameterNames.ID_TOKEN);
    LoginResponse loginResponse =
        LoginResponse.builder()
            .accessToken(tokenResponse.getAccessToken().getTokenValue())
            .idToken(idToken)
            .refreshToken(
                tokenResponse.getRefreshToken() != null
                    ? tokenResponse.getRefreshToken().getTokenValue()
                    : null)
            .accessTokenType(tokenResponse.getAccessToken().getTokenType().getValue())
            .expiresIn(getExpiresIn(tokenResponse))
            .build();
    return JsonUtil.convertValue(R.ok(loginResponse), Map.class);
  }
}
