package top.wecoding.iam.framework.cache;

import org.jetbrains.annotations.NotNull;
import top.wecoding.commons.core.cache.CacheKeyBuilder;
import top.wecoding.iam.common.constant.RedisConstant;

import java.time.Duration;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
public class AuthorizationConsentCacheKeyBuilder implements CacheKeyBuilder {

  @NotNull
  @Override
  public String getPrefix() {
    return RedisConstant.AUTHORIZATION_CONSENT;
  }

  @Override
  public Duration getExpire() {
    return Duration.ofMinutes(10L);
  }
}
