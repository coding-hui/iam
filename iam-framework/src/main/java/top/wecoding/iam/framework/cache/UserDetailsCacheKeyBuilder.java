package top.wecoding.iam.framework.cache;

import org.springframework.lang.NonNull;
import top.wecoding.commons.core.cache.CacheKeyBuilder;
import top.wecoding.iam.common.constant.RedisConstant;

/**
 * @author liuyuhui
 * @date 2022/10/8
 */
public class UserDetailsCacheKeyBuilder implements CacheKeyBuilder {

  @NonNull
  @Override
  public String getPrefix() {
    return RedisConstant.USER_DETAILS;
  }
}
