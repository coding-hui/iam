package top.wecoding.iam.server.cache;

import top.wecoding.core.cache.CacheKeyBuilder;
import top.wecoding.iam.common.constant.RedisConstant;

/**
 * @author liuyuhui
 * @date 2022/10/1
 * @qq 1515418211
 */
public class UserFailCountCacheKeyBuilder implements CacheKeyBuilder {

  @Override
  public String getPrefix() {
    return RedisConstant.USER_FAIL_COUNT;
  }
}
