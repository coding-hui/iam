package top.wecoding.iam.server.security.authorization.token;

import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenClaimsContext;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenClaimsSet;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenCustomizer;
import top.wecoding.iam.common.constant.SecurityConstants;
import top.wecoding.iam.common.userdetails.LoginUser;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
public class WecodingOAuth2TokenCustomizer
    implements OAuth2TokenCustomizer<OAuth2TokenClaimsContext> {

  /**
   * Customize the OAuth 2.0 Token attributes.
   *
   * @param context the context containing the OAuth 2.0 Token attributes
   */
  @Override
  public void customize(OAuth2TokenClaimsContext context) {
    OAuth2TokenClaimsSet.Builder claims = context.getClaims();
    claims.claim(SecurityConstants.DETAILS_LICENSE, SecurityConstants.PROJECT_LICENSE);
    claims.claim(SecurityConstants.CLIENT_ID, context.getAuthorizationGrant().getName());
    // 客户端模式不返回具体用户信息
    if (SecurityConstants.CLIENT_CREDENTIALS.equals(
        context.getAuthorizationGrantType().getValue())) {
      return;
    }
    LoginUser loginUser = (LoginUser) context.getPrincipal().getPrincipal();
    claims.claim(SecurityConstants.USER_ID, loginUser.userInfo().getUserId());
  }
}