package top.wecoding.iam.server.authentication.password;

import java.util.Map;
import java.util.Set;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.core.AuthorizationGrantType;
import top.wecoding.iam.server.authentication.base.OAuth2ResourceOwnerBaseAuthenticationToken;

/**
 * 密码授权token信息
 *
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
public class OAuth2ResourceOwnerPasswordAuthenticationToken
    extends OAuth2ResourceOwnerBaseAuthenticationToken {

  public OAuth2ResourceOwnerPasswordAuthenticationToken(
      AuthorizationGrantType authorizationGrantType,
      Authentication clientPrincipal,
      Set<String> scopes,
      Map<String, Object> additionalParameters) {
    super(authorizationGrantType, clientPrincipal, scopes, additionalParameters);
  }
}
