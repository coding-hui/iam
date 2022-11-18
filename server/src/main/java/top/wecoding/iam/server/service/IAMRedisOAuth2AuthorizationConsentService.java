package top.wecoding.iam.server.service;

import lombok.RequiredArgsConstructor;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationConsent;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationConsentService;
import org.springframework.stereotype.Service;
import org.springframework.util.Assert;
import top.wecoding.iam.framework.cache.AuthorizationConsentCacheKeyBuilder;

import java.util.concurrent.TimeUnit;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Service
@RequiredArgsConstructor
public class IAMRedisOAuth2AuthorizationConsentService
    implements OAuth2AuthorizationConsentService {

  private static final Long TIMEOUT = 10L;

  private final RedisTemplate<String, Object> redisTemplate;

  @Override
  public void save(OAuth2AuthorizationConsent authorizationConsent) {
    Assert.notNull(authorizationConsent, "authorizationConsent cannot be null");
    redisTemplate
        .opsForValue()
        .set(buildKey(authorizationConsent), authorizationConsent, TIMEOUT, TimeUnit.MINUTES);
  }

  @Override
  public void remove(OAuth2AuthorizationConsent authorizationConsent) {
    Assert.notNull(authorizationConsent, "authorizationConsent cannot be null");
    redisTemplate.delete(buildKey(authorizationConsent));
  }

  @Override
  public OAuth2AuthorizationConsent findById(String registeredClientId, String principalName) {
    Assert.hasText(registeredClientId, "registeredClientId cannot be empty");
    Assert.hasText(principalName, "principalName cannot be empty");
    return (OAuth2AuthorizationConsent)
        redisTemplate.opsForValue().get(buildKey(registeredClientId, principalName));
  }

  private String buildKey(OAuth2AuthorizationConsent authorizationConsent) {
    return buildKey(
        authorizationConsent.getRegisteredClientId(), authorizationConsent.getPrincipalName());
  }

  private String buildKey(String registeredClientId, String principalName) {
    return new AuthorizationConsentCacheKeyBuilder()
        .key(registeredClientId, principalName)
        .getKey();
  }
}
