package top.wecoding.iam.common.cache;

import top.wecoding.core.cache.CacheKeyBuilder;
import top.wecoding.iam.common.constant.RedisConstant;

import java.time.Duration;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
public class AuthorizationConsentCacheKeyBuilder implements CacheKeyBuilder {

  @Override
  public String getPrefix() {
    return RedisConstant.AUTHORIZATION_CONSENT;
  }

  @Override
  public Duration getExpire() {
    return Duration.ofMinutes(10L);
  }
}
