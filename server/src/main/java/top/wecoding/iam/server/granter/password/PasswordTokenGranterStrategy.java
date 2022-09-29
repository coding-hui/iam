package top.wecoding.iam.server.granter.password;

import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import top.wecoding.iam.common.model.LoginUser;
import top.wecoding.iam.common.model.request.CreateUserRequest;
import top.wecoding.iam.common.model.request.TokenRequest;
import top.wecoding.iam.server.granter.AbstractTokenGranterStrategy;
import top.wecoding.iam.server.mapper.UserMapper;
import top.wecoding.iam.server.pojo.User;

import javax.annotation.Resource;

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

  @Resource UserMapper userMapper;

  @Override
  protected LoginUser loadUserInfo(TokenRequest tokenRequest) {
    String account = tokenRequest.getAccount();
    String password = tokenRequest.getPassword();

    CreateUserRequest createUserRequest =
        CreateUserRequest.builder().username(account).password(password).build();

    User user = userMapper.getInfoByUsernameAntTenantName(createUserRequest.getUsername(), "");

    return new LoginUser();
  }
}
