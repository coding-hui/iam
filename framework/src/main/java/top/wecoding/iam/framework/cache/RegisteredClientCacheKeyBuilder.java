package top.wecoding.iam.framework.cache;

import top.wecoding.core.cache.CacheKeyBuilder;
import top.wecoding.iam.common.constant.RedisConstant;

import java.time.Duration;

/**
 * @author liuyuhui
 * @date 2022/10/16
 * @qq 1515418211
 */
public class RegisteredClientCacheKeyBuilder implements CacheKeyBuilder {

  @Override
  public String getPrefix() {
    return RedisConstant.CLIENT_DETAILS_KEY;
  }

  @Override
  public Duration getExpire() {
    return Duration.ofMinutes(30L);
  }
}
