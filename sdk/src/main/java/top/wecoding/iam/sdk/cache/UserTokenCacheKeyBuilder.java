package top.wecoding.iam.sdk.cache;

import java.time.Duration;
import org.springframework.lang.NonNull;
import top.wecoding.core.cache.CacheKeyBuilder;
import top.wecoding.iam.common.constant.RedisConstant;

/**
 * @author liuyuhui
 * @date 2022/9/11
 * @qq 1515418211
 */
public class UserTokenCacheKeyBuilder implements CacheKeyBuilder {

  @NonNull
  @Override
  public String getPrefix() {
    return RedisConstant.USER_TOKEN;
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
