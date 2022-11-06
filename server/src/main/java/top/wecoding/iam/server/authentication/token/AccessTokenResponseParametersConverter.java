package top.wecoding.iam.server.authentication.token;

import cn.hutool.core.bean.BeanUtil;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.HashMap;
import java.util.Map;
import org.springframework.core.convert.converter.Converter;
import org.springframework.security.oauth2.core.endpoint.OAuth2AccessTokenResponse;
import org.springframework.security.oauth2.core.endpoint.OAuth2ParameterNames;
import org.springframework.stereotype.Component;
import org.springframework.util.CollectionUtils;
import org.springframework.util.StringUtils;
import top.wecoding.core.result.R;

/**
 * @author liuyuhui
 * @date 2022/11/5
 * @qq 1515418211
 */
@Component
public class AccessTokenResponseParametersConverter
    implements Converter<OAuth2AccessTokenResponse, Map<String, Object>> {

  private static long getExpiresIn(OAuth2AccessTokenResponse tokenResponse) {
    if (tokenResponse.getAccessToken().getExpiresAt() != null) {
      return ChronoUnit.SECONDS.between(
          Instant.now(), tokenResponse.getAccessToken().getExpiresAt());
    }
    return -1;
  }

  @Override
  public Map<String, Object> convert(OAuth2AccessTokenResponse tokenResponse) {
    Map<String, Object> parameters = new HashMap<>();
    parameters.put(
        OAuth2ParameterNames.ACCESS_TOKEN, tokenResponse.getAccessToken().getTokenValue());
    parameters.put(
        OAuth2ParameterNames.TOKEN_TYPE, tokenResponse.getAccessToken().getTokenType().getValue());
    parameters.put(OAuth2ParameterNames.EXPIRES_IN, getExpiresIn(tokenResponse));
    if (!CollectionUtils.isEmpty(tokenResponse.getAccessToken().getScopes())) {
      parameters.put(
          OAuth2ParameterNames.SCOPE,
          StringUtils.collectionToDelimitedString(tokenResponse.getAccessToken().getScopes(), " "));
    }
    if (tokenResponse.getRefreshToken() != null) {
      parameters.put(
          OAuth2ParameterNames.REFRESH_TOKEN, tokenResponse.getRefreshToken().getTokenValue());
    }
    if (!CollectionUtils.isEmpty(tokenResponse.getAdditionalParameters())) {
      for (Map.Entry<String, Object> entry : tokenResponse.getAdditionalParameters().entrySet()) {
        parameters.put(entry.getKey(), entry.getValue());
      }
    }

    return BeanUtil.beanToMap(R.ok(parameters));
  }
}
