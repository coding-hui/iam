package top.wecoding.iam.server.security.authorization.authentication.password;

import org.springframework.security.core.Authentication;
import top.wecoding.iam.common.enums.AuthType;
import top.wecoding.iam.common.model.request.LoginRequest;
import top.wecoding.iam.server.security.authorization.authentication.OAuth2ResourceOwnerBaseAuthenticationToken;

import java.util.Set;

/**
 * 密码授权token信息
 *
 * @author liuyuhui
 * @date 2022/10/3
 */
public class OAuth2ResourceOwnerPasswordAuthenticationToken
    extends OAuth2ResourceOwnerBaseAuthenticationToken {

  public OAuth2ResourceOwnerPasswordAuthenticationToken(
      AuthType authType,
      Authentication clientPrincipal,
      Set<String> scopes,
      LoginRequest loginRequest) {
    super(authType, clientPrincipal, scopes, loginRequest);
  }
}
