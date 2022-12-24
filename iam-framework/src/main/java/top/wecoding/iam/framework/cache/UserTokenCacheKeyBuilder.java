package top.wecoding.iam.framework.cache;

import java.time.Duration;
import org.jetbrains.annotations.NotNull;
import org.springframework.lang.NonNull;
import top.wecoding.commons.core.cache.CacheKeyBuilder;
import top.wecoding.iam.common.constant.RedisConstant;

/**
 * @author liuyuhui
 * @date 2022/9/11
 */
public class UserTokenCacheKeyBuilder implements CacheKeyBuilder {

  @NotNull
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
