package top.wecoding.iam.granter.refresh;

import cn.hutool.core.util.StrUtil;
import io.jsonwebtoken.Claims;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import top.wecoding.auth.model.LoginUser;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.core.constant.TokenConstant;
import top.wecoding.core.exception.user.UnauthorizedException;
import top.wecoding.iam.granter.AbstractTokenGranterStrategy;
import top.wecoding.iam.model.request.TokenRequest;
import top.wecoding.iam.service.UserService;
import top.wecoding.jwt.util.JwtUtils;

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

    if (JwtUtils.isTokenExpired(refreshToken)) {
      throw new UnauthorizedException();
    }

    Claims claims = JwtUtils.parseToken(refreshToken);
    String tokenType = JwtUtils.getValue(claims, TokenConstant.TOKEN_TYPE);
    String account = JwtUtils.getValue(claims, SecurityConstants.DETAILS_ACCOUNT);
    if (!TokenConstant.REFRESH_TOKEN.equals(tokenType)) {
      throw new UnauthorizedException();
    }

    return new LoginUser();
  }
}
