package top.wecoding.iam.granter;

import cn.hutool.core.util.StrUtil;
import cn.hutool.crypto.digest.BCrypt;
import cn.hutool.extra.spring.SpringUtil;
import lombok.extern.slf4j.Slf4j;
import top.wecoding.auth.model.AuthInfo;
import top.wecoding.auth.model.LoginUser;
import top.wecoding.core.enums.rest.CommonErrorCodeEnum;
import top.wecoding.core.exception.user.UnauthorizedException;
import top.wecoding.core.util.AssertUtils;
import top.wecoding.iam.model.request.TokenRequest;
import top.wecoding.iam.service.UserService;
import top.wecoding.jwt.model.TokenInfo;
import top.wecoding.jwt.util.JwtUtils;
import top.wecoding.security.helper.AuthHelper;
import top.wecoding.security.provider.ClientDetails;
import top.wecoding.security.provider.ClientDetailsService;

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
  protected static final ClientDetailsService clientDetailsService;

  static {
    userService = SpringUtil.getBean(UserService.class);
    clientDetailsService = SpringUtil.getBean(ClientDetailsService.class);
  }

  @Override
  public final AuthInfo grant(TokenRequest tokenRequest) {
    String[] tokens = determineClient();
    String clientId = tokens[0];
    String clientSecret = tokens[1];

    // 获取客户端
    ClientDetails clientDetails = loadClientDetails(clientId);

    // 进行客户端检验
    if (!validateClient(clientDetails, clientId, clientSecret)) {
      throw new UnauthorizedException();
    }

    // 具体策略登录
    LoginUser loginUser = loadUserInfo(tokenRequest);

    // 登录后构建 Token
    AuthInfo authInfo = createToken(loginUser, clientDetails);

    // 更新登录信息
    afterLoginSuccess(loginUser);

    return authInfo;
  }

  protected String[] determineClient() {
    String[] tokens = JwtUtils.extractAndDecodeHeader();
    AssertUtils.isTrue(tokens.length == 2, CommonErrorCodeEnum.COMMON_ERROR, "client vaild error");
    return tokens;
  }

  /** 获取客户配置信息 */
  protected ClientDetails loadClientDetails(String clientId) {
    // 获取客户端信息
    return clientDetailsService.loadClientByClientId(clientId);
  }

  /** 校验Client */
  protected boolean validateClient(
      ClientDetails clientDetails, String clientId, String clientSecret) {
    if (clientDetails != null) {
      return StrUtil.equals(clientId, clientDetails.getClientId())
          && BCrypt.checkpw(clientSecret, clientDetails.getClientSecret());
    }
    return false;
  }

  /** 创建令牌 */
  protected AuthInfo createToken(LoginUser loginUser, ClientDetails clientDetails) {
    loginUser.setClientId(clientDetails.getClientId());
    return AuthHelper.ofAuthInfo(loginUser, new TokenInfo());
  }

  /** 登录成功之后 */
  protected void afterLoginSuccess(LoginUser loginUser) {
    // 登录日志记录
  }

  /** 登录获取用户信息 */
  protected abstract LoginUser loadUserInfo(TokenRequest tokenRequest);
}
