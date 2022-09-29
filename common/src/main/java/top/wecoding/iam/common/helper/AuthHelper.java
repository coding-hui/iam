package top.wecoding.iam.common.helper;

import cn.hutool.core.util.IdUtil;
import cn.hutool.extra.spring.SpringUtil;
import lombok.experimental.UtilityClass;
import org.springframework.data.redis.connection.RedisStringCommands;
import org.springframework.data.redis.core.RedisCallback;
import org.springframework.data.redis.core.types.Expiration;
import org.springframework.util.Base64Utils;
import top.wecoding.core.cache.CacheKey;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.core.constant.TokenConstant;
import top.wecoding.core.util.ProtostuffUtil;
import top.wecoding.iam.common.cache.UserTokenCacheKeyBuilder;
import top.wecoding.iam.common.constant.UserConstant;
import top.wecoding.iam.common.model.AuthInfo;
import top.wecoding.iam.common.model.LoginUser;
import top.wecoding.iam.common.util.TokenUtils;
import top.wecoding.jwt.model.JwtPayLoad;
import top.wecoding.jwt.model.TokenInfo;
import top.wecoding.jwt.props.JwtProperties;
import top.wecoding.redis.repository.base.CacheOperatorPlus;

/**
 * @author liuyuhui
 * @date 2022/9/11
 * @qq 1515418211
 */
@UtilityClass
public class AuthHelper {

  private static final JwtProperties jwtProperties;
  private static final CacheOperatorPlus cacheOperatorPlus;

  static {
    jwtProperties = SpringUtil.getBean(JwtProperties.class);
    cacheOperatorPlus = SpringUtil.getBean(CacheOperatorPlus.class);
  }

  public void invalidWebSession(String userId) {
    cacheOperatorPlus.del(userId);
  }

  public void setWebSession(AuthInfo webSession) {
    byte[] bytes = Base64Utils.decodeFromString(webSession.getAccessToken());
    String[] split = new String(bytes).split(UserConstant.SPACE);

    CacheKey key = new UserTokenCacheKeyBuilder().key(split[1]);
    cacheOperatorPlus
        .getRedisTemplate()
        .executePipelined(
            (RedisCallback<?>)
                connection -> {
                  assert key.getExpire() != null;
                  connection.set(
                      (webSession.getTenantId() + UserConstant.DASH + webSession.getUserId())
                          .getBytes(),
                      key.getKey().getBytes(),
                      Expiration.seconds(key.getExpire().getSeconds()),
                      RedisStringCommands.SetOption.ifAbsent());
                  connection.set(
                      key.getKey().getBytes(),
                      ProtostuffUtil.serialize(webSession),
                      Expiration.seconds(key.getExpire().getSeconds()),
                      RedisStringCommands.SetOption.ifAbsent());
                  return null;
                });
  }

  public AuthInfo getWebSession(String token) {
    CacheKey key = new UserTokenCacheKeyBuilder().key(token);
    AuthInfo session = cacheOperatorPlus.get(key);
    if (session != null && session.getUserId() != null) {
      cacheOperatorPlus.expire(key);
    }
    return session;
  }

  public AuthInfo ofAuthInfo(LoginUser loginUser, TokenInfo tokenInfo) {
    loginUser.setUuid(IdUtil.fastUUID());
    JwtPayLoad jwtPayLoad =
        JwtPayLoad.builder()
            .uuid(loginUser.getUuid())
            .userId(loginUser.getUserId())
            .account(loginUser.getAccount())
            .clientId(loginUser.getClientId())
            .realName(loginUser.getRealName())
            .build();

    return AuthInfo.builder()
        .accessToken(tokenInfo.getToken())
        .expireMillis(tokenInfo.getExpiresIn())
        .expiration(tokenInfo.getExpiration())
        .tokenType(TokenConstant.ACCESS_TOKEN)
        .refreshToken(TokenUtils.createRefreshToken(jwtPayLoad).getToken())
        .uuid(jwtPayLoad.getUuid())
        .userId(jwtPayLoad.getUserId())
        .account(jwtPayLoad.getAccount())
        .realName(jwtPayLoad.getRealName())
        .clientId(jwtPayLoad.getClientId())
        .license(SecurityConstants.PROJECT_LICENSE)
        .build();
  }
}
