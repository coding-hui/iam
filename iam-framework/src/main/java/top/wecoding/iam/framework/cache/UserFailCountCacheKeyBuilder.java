package top.wecoding.iam.framework.cache;

import org.jetbrains.annotations.NotNull;
import top.wecoding.commons.core.cache.CacheKeyBuilder;
import top.wecoding.iam.common.constant.RedisConstant;

/**
 * @author liuyuhui
 * @date 2022/10/1
 */
public class UserFailCountCacheKeyBuilder implements CacheKeyBuilder {

  @NotNull
  @Override
  public String getPrefix() {
    return RedisConstant.USER_FAIL_COUNT;
  }
}