package top.wecoding.iam.sdk.service;

import static org.springframework.security.oauth2.core.endpoint.OAuth2ParameterNames.*;

import java.time.temporal.ChronoUnit;
import java.util.ArrayList;
import java.util.List;
import java.util.Objects;
import java.util.concurrent.TimeUnit;
import lombok.RequiredArgsConstructor;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.data.redis.serializer.RedisSerializer;
import org.springframework.lang.Nullable;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.security.oauth2.core.OAuth2AuthorizationCode;
import org.springframework.security.oauth2.core.OAuth2RefreshToken;
import org.springframework.security.oauth2.core.OAuth2TokenType;
import org.springframework.security.oauth2.server.authorization.OAuth2Authorization;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.util.Assert;
import top.wecoding.iam.common.cache.UserTokenCacheKeyBuilder;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@RequiredArgsConstructor
public class IAMRedisOAuth2AuthorizationService implements OAuth2AuthorizationService {

  private static final Long TIMEOUT = 10L;

  private final RedisTemplate<String, Object> redisTemplate;

  @Override
  public void save(OAuth2Authorization authorization) {
    Assert.notNull(authorization, "authorization cannot be null");

    if (isState(authorization)) {
      String token = authorization.getAttribute(STATE);
      redisTemplate.setValueSerializer(RedisSerializer.java());
      redisTemplate
          .opsForValue()
          .set(buildKey(STATE, token), authorization, TIMEOUT, TimeUnit.MINUTES);
    }

    if (isCode(authorization)) {
      OAuth2Authorization.Token<OAuth2AuthorizationCode> authorizationCode =
          authorization.getToken(OAuth2AuthorizationCode.class);
      OAuth2AuthorizationCode authorizationCodeToken = authorizationCode.getToken();
      long between =
          ChronoUnit.MINUTES.between(
              authorizationCodeToken.getIssuedAt(), authorizationCodeToken.getExpiresAt());
      redisTemplate.setValueSerializer(RedisSerializer.java());
      redisTemplate
          .opsForValue()
          .set(
              buildKey(CODE, authorizationCodeToken.getTokenValue()),
              authorization,
              between,
              TimeUnit.MINUTES);
    }

    if (isRefreshToken(authorization)) {
      OAuth2RefreshToken refreshToken = authorization.getRefreshToken().getToken();
      long between =
          ChronoUnit.SECONDS.between(refreshToken.getIssuedAt(), refreshToken.getExpiresAt());
      redisTemplate.setValueSerializer(RedisSerializer.java());
      redisTemplate
          .opsForValue()
          .set(
              buildKey(REFRESH_TOKEN, refreshToken.getTokenValue()),
              authorization,
              between,
              TimeUnit.SECONDS);
    }

    if (isAccessToken(authorization)) {
      OAuth2AccessToken accessToken = authorization.getAccessToken().getToken();
      long between =
          ChronoUnit.SECONDS.between(accessToken.getIssuedAt(), accessToken.getExpiresAt());
      redisTemplate.setValueSerializer(RedisSerializer.java());
      redisTemplate
          .opsForValue()
          .set(
              buildKey(ACCESS_TOKEN, accessToken.getTokenValue()),
              authorization,
              between,
              TimeUnit.SECONDS);
    }
  }

  @Override
  public void remove(OAuth2Authorization authorization) {
    Assert.notNull(authorization, "authorization cannot be null");

    List<String> keys = new ArrayList<>();
    if (isState(authorization)) {
      String token = authorization.getAttribute(STATE);
      keys.add(buildKey(STATE, token));
    }

    if (isCode(authorization)) {
      OAuth2Authorization.Token<OAuth2AuthorizationCode> authorizationCode =
          authorization.getToken(OAuth2AuthorizationCode.class);
      OAuth2AuthorizationCode authorizationCodeToken = authorizationCode.getToken();
      keys.add(buildKey(CODE, authorizationCodeToken.getTokenValue()));
    }

    if (isRefreshToken(authorization)) {
      OAuth2RefreshToken refreshToken = authorization.getRefreshToken().getToken();
      keys.add(buildKey(REFRESH_TOKEN, refreshToken.getTokenValue()));
    }

    if (isAccessToken(authorization)) {
      OAuth2AccessToken accessToken = authorization.getAccessToken().getToken();
      keys.add(buildKey(ACCESS_TOKEN, accessToken.getTokenValue()));
    }
    redisTemplate.delete(keys);
  }

  @Override
  @Nullable
  public OAuth2Authorization findById(String id) {
    throw new UnsupportedOperationException();
  }

  @Override
  @Nullable
  public OAuth2Authorization findByToken(String token, @Nullable OAuth2TokenType tokenType) {
    Assert.hasText(token, "token cannot be empty");
    Assert.notNull(tokenType, "tokenType cannot be empty");
    redisTemplate.setValueSerializer(RedisSerializer.java());
    return (OAuth2Authorization)
        redisTemplate.opsForValue().get(buildKey(tokenType.getValue(), token));
  }

  private String buildKey(String type, String id) {
    return new UserTokenCacheKeyBuilder().key(type, id).getKey();
  }

  private boolean isState(OAuth2Authorization authorization) {
    return Objects.nonNull(authorization.getAttribute(STATE));
  }

  private boolean isCode(OAuth2Authorization authorization) {
    OAuth2Authorization.Token<OAuth2AuthorizationCode> authorizationCode =
        authorization.getToken(OAuth2AuthorizationCode.class);
    return Objects.nonNull(authorizationCode);
  }

  private boolean isRefreshToken(OAuth2Authorization authorization) {
    return Objects.nonNull(authorization.getRefreshToken());
  }

  private boolean isAccessToken(OAuth2Authorization authorization) {
    return Objects.nonNull(authorization.getAccessToken());
  }
}
