package top.wecoding.iam.server.util.cache;

import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import top.wecoding.core.cache.CacheKey;
import top.wecoding.iam.server.cache.UserFailCountCacheKeyBuilder;

/**
 * @author liuyuhui
 * @date 2022/10/2
 * @qq 1515418211
 */
public class UserFailCountCacheKeyBuilderTest {

  @Test
  void genKey() {
    CacheKey key = new UserFailCountCacheKeyBuilder().key("1");
    Assertions.assertNotNull(key);
    Assertions.assertNotNull(key.getExpire());
  }
}
