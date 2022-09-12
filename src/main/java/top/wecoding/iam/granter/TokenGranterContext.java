package top.wecoding.iam.granter;

import cn.hutool.extra.spring.SpringUtil;
import lombok.AllArgsConstructor;
import top.wecoding.auth.model.AuthInfo;
import top.wecoding.core.exception.user.UnauthorizedException;
import top.wecoding.iam.granter.password.PasswordTokenGranterStrategy;
import top.wecoding.iam.granter.refresh.RefreshTokenGranterStrategy;
import top.wecoding.iam.model.request.TokenRequest;

import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @author liuyuhui
 * @date 2022
 * @qq 1515418211
 */
@AllArgsConstructor
public class TokenGranterContext {

  private static final Map<String, TokenGranterStrategy> GRANTER_STRATEGY_MAP =
      new ConcurrentHashMap<>();

  static {
    GRANTER_STRATEGY_MAP.put(
        PasswordTokenGranterStrategy.GRANT_TYPE,
        SpringUtil.getBean(PasswordTokenGranterStrategy.class));
    GRANTER_STRATEGY_MAP.put(
        RefreshTokenGranterStrategy.GRANT_TYPE,
        SpringUtil.getBean(RefreshTokenGranterStrategy.class));
  }

  public static AuthInfo grant(TokenRequest tokenRequest) {
    String grantType = tokenRequest.getGrantType();
    if (!GRANTER_STRATEGY_MAP.containsKey(grantType)) {
      throw new UnauthorizedException();
    }

    TokenGranterStrategy tokenGranterStrategy = GRANTER_STRATEGY_MAP.get(grantType);

    return tokenGranterStrategy.grant(tokenRequest);
  }
}
