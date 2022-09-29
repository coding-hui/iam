package top.wecoding.iam.common.cache;

import org.springframework.lang.NonNull;
import top.wecoding.core.cache.CacheKeyBuilder;
import top.wecoding.iam.common.constant.CacheKeyDefinition;

import java.time.Duration;

/**
 * @author liuyuhui
 * @date 2022/9/11
 * @qq 1515418211
 */
public class UserTokenCacheKeyBuilder implements CacheKeyBuilder {

  @NonNull
  @Override
  public String getPrefix() {
    return CacheKeyDefinition.LOGIN_USER_INFO;
  }

  @Override
  public String getTenant() {
    return null;
  }

  @Override
  public Duration getExpire() {
    return Duration.ofMinutes(30);
  }
}
