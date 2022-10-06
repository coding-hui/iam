package top.wecoding.iam.server.cache;

import top.wecoding.core.cache.CacheKeyBuilder;
import top.wecoding.iam.common.constant.RedisConstant;

import java.time.Duration;

/**
 * @author liuyuhui
 * @date 2022/10/1
 * @qq 1515418211
 */
public class UserFailCountCacheKeyBuilder implements CacheKeyBuilder {

  @Override
  public Duration getExpire() {
    return Duration.ofHours(1L);
  }

  @Override
  public String getPrefix() {
    return RedisConstant.USER_FAIL_COUNT;
  }
}
