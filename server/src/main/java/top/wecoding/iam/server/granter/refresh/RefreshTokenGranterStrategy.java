package top.wecoding.iam.server.granter.refresh;

import cn.hutool.core.util.StrUtil;
import io.jsonwebtoken.Claims;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.core.constant.TokenConstant;
import top.wecoding.core.exception.user.UnauthorizedException;
import top.wecoding.iam.common.model.LoginUser;
import top.wecoding.iam.common.model.request.TokenRequest;
import top.wecoding.iam.server.granter.AbstractTokenGranterStrategy;
import top.wecoding.iam.server.service.UserService;
import top.wecoding.jwt.helper.JWTHelper;

/**
 * 刷新令牌换取访问令牌
 *
 * @author liuyuhui
 * @date 2022/5/4
 * @qq 1515418211
 */
@Slf4j
@Component
@AllArgsConstructor
public class RefreshTokenGranterStrategy extends AbstractTokenGranterStrategy {

  public static final String GRANT_TYPE = "refresh_token";

  private final UserService userService;

  @Override
  protected LoginUser loadUserInfo(TokenRequest tokenRequest) {
    String grantType = tokenRequest.getGrantType();
    String refreshToken = tokenRequest.getRefreshToken();

    if (StrUtil.hasBlank(grantType, refreshToken) || !GRANT_TYPE.equals(grantType)) {
      throw new UnauthorizedException();
    }

    if (JWTHelper.isTokenExpired(refreshToken)) {
      throw new UnauthorizedException();
    }

    Claims claims = JWTHelper.parseToken(refreshToken);
    String tokenType = JWTHelper.getValue(claims, TokenConstant.TOKEN_TYPE);
    String account = JWTHelper.getValue(claims, SecurityConstants.DETAILS_ACCOUNT);
    if (!TokenConstant.REFRESH_TOKEN.equals(tokenType)) {
      throw new UnauthorizedException();
    }

    return new LoginUser();
  }
}