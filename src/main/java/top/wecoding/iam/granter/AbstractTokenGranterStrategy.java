package top.wecoding.iam.granter;

import cn.hutool.core.lang.UUID;
import cn.hutool.extra.spring.SpringUtil;
import lombok.extern.slf4j.Slf4j;
import top.wecoding.auth.model.AuthInfo;
import top.wecoding.auth.model.LoginUser;
import top.wecoding.core.enums.rest.CommonErrorCodeEnum;
import top.wecoding.core.exception.user.UnauthorizedException;
import top.wecoding.core.util.AssertUtils;
import top.wecoding.iam.model.request.TokenRequest;
import top.wecoding.iam.service.UserService;
import top.wecoding.jwt.helper.JWTHelper;
import top.wecoding.jwt.model.JwtPayLoad;
import top.wecoding.security.helper.AuthHelper;
import top.wecoding.security.helper.SecurityHelper;
import top.wecoding.security.provider.ClientDetails;

/**
 * 抽象登录父类，封装一些公共方法，检校逻辑
 *
 * @author liuyuhui
 * @date 2022
 * @qq 1515418211
 */
@Slf4j
public abstract class AbstractTokenGranterStrategy implements TokenGranterStrategy {

  protected static final UserService userService;

  static {
    userService = SpringUtil.getBean(UserService.class);
  }

  @Override
  public final AuthInfo grant(TokenRequest tokenRequest) {
    String[] tokens = determineClient();
    String clientId = tokens[0];
    String clientSecret = tokens[1];

    ClientDetails clientDetails = loadClientDetails(clientId);

    if (!validateClient(clientDetails, clientId, clientSecret)) {
      throw new UnauthorizedException();
    }

    LoginUser loginUser = loadUserInfo(tokenRequest);

    AuthInfo authInfo = createToken(loginUser, clientDetails);

    afterLoginSuccess(loginUser);

    return authInfo;
  }

  protected String[] determineClient() {
    String[] tokens = SecurityHelper.extractAndDecodeHeader();
    AssertUtils.isTrue(
        tokens.length == 2,
        CommonErrorCodeEnum.COMMON_ERROR,
        "Failed to determine client authentication information in request header");
    return tokens;
  }

  protected ClientDetails loadClientDetails(String clientId) {
    return SecurityHelper.clientDetails(clientId);
  }

  protected boolean validateClient(
      ClientDetails clientDetails, String clientId, String clientSecret) {
    return SecurityHelper.validateClient(clientDetails, clientId, clientSecret);
  }

  protected abstract LoginUser loadUserInfo(TokenRequest tokenRequest);

  protected AuthInfo createToken(LoginUser loginUser, ClientDetails clientDetails) {
    loginUser.setUuid(UUID.randomUUID().toString());
    loginUser.setClientId(clientDetails.getClientId());

    JwtPayLoad payload =
        JwtPayLoad.builder()
            .uuid(loginUser.getUuid())
            .userId(loginUser.getUserId())
            .account(loginUser.getAccount())
            .clientId(loginUser.getClientId())
            .build();

    AuthInfo authInfo =
        AuthHelper.ofAuthInfo(
            loginUser,
            JWTHelper.createToken(payload, clientDetails.getAccessTokenValiditySeconds()));

    AuthHelper.setWebSession(authInfo);

    return authInfo;
  }

  protected void afterLoginSuccess(LoginUser loginUser) {
    log.info("user: {} login success.", loginUser.getAccount());
  }
}
