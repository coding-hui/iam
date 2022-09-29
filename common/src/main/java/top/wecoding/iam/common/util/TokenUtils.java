package top.wecoding.iam.common.util;

import cn.hutool.core.convert.Convert;
import cn.hutool.extra.spring.SpringUtil;
import org.springframework.util.Base64Utils;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.core.constant.TokenConstant;
import top.wecoding.iam.common.constant.UserConstant;
import top.wecoding.iam.common.enums.AuthType;
import top.wecoding.jwt.helper.JWTHelper;
import top.wecoding.jwt.model.JwtPayLoad;
import top.wecoding.jwt.model.TokenInfo;
import top.wecoding.jwt.props.JwtProperties;
import top.wecoding.redis.repository.base.CacheOperatorPlus;

import java.util.HashMap;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.ThreadLocalRandom;

/**
 * @author liuyuhui
 * @date 2022/9/11
 * @qq 1515418211
 */
public class TokenUtils {

  private static final JwtProperties jwtProperties;
  private static final CacheOperatorPlus cacheOperatorPlus;

  static {
    jwtProperties = SpringUtil.getBean(JwtProperties.class);
    cacheOperatorPlus = SpringUtil.getBean(CacheOperatorPlus.class);
  }

  public static String createWebToken() {
    ThreadLocalRandom random = ThreadLocalRandom.current();
    String uuid =
        new UUID(random.nextLong(), random.nextLong())
            .toString()
            .replace(UserConstant.DASH, UserConstant.EMPTY);
    String token = AuthType.WEB.code() + UserConstant.SPACE + uuid;
    return Base64Utils.encodeToString(token.getBytes());
  }

  public static String createApiToken() {
    ThreadLocalRandom random = ThreadLocalRandom.current();
    String uuid =
        new UUID(random.nextLong(), random.nextLong())
            .toString()
            .replace(UserConstant.DASH, UserConstant.EMPTY);
    String token = AuthType.API_TOKEN.code() + UserConstant.SPACE + uuid;
    return Base64Utils.encodeToString(token.getBytes());
  }

  public static TokenInfo createRefreshToken(JwtPayLoad jwtPayLoad) {
    Map<String, String> claims = new HashMap<>(16);
    claims.put(TokenConstant.TOKEN_TYPE, TokenConstant.REFRESH_TOKEN);
    claims.put(SecurityConstants.USER_KEY, jwtPayLoad.getUuid());
    claims.put(SecurityConstants.DETAILS_USER_ID, Convert.toStr(jwtPayLoad.getUserId()));
    claims.put(SecurityConstants.DETAILS_ACCOUNT, jwtPayLoad.getAccount());
    claims.put(SecurityConstants.DETAILS_CLIENT_ID, jwtPayLoad.getClientId());
    return JWTHelper.createToken(claims, jwtProperties.getRefreshExpire());
  }
}
