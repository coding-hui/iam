package top.wecoding.iam.granter.password;

import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import top.wecoding.auth.model.LoginUser;
import top.wecoding.iam.granter.AbstractTokenGranterStrategy;
import top.wecoding.iam.model.request.TokenRequest;

/**
 * 密码登录
 *
 * @author liuyuhui
 * @date 2022
 * @qq 1515418211
 */
@Slf4j
@Component
@AllArgsConstructor
public class PasswordTokenGranterStrategy extends AbstractTokenGranterStrategy {

  public static final String GRANT_TYPE = "password";

  @Override
  protected LoginUser loadUserInfo(TokenRequest tokenRequest) {
    String account = tokenRequest.getAccount();
    String password = tokenRequest.getPassword();
    return new LoginUser();
  }
}
