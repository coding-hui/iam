package top.wecoding.iam.server.security.authorization.token;

import org.springframework.security.oauth2.jwt.JwtClaimsSet;
import org.springframework.security.oauth2.server.authorization.token.JwtEncodingContext;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenCustomizer;
import top.wecoding.iam.common.constant.SecurityConstants;
import top.wecoding.iam.common.userdetails.LoginUser;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
public class WecodingOAuth2JwtTokenCustomizer implements OAuth2TokenCustomizer<JwtEncodingContext> {

  /**
   * Customize the OAuth 2.0 Token attributes.
   *
   * @param context the context containing the OAuth 2.0 Token attributes
   */
  @Override
  public void customize(JwtEncodingContext context) {
    JwtClaimsSet.Builder claims = context.getClaims();
    // 客户端模式不返回具体用户信息
    if (SecurityConstants.CLIENT_CREDENTIALS.equals(
        context.getAuthorizationGrantType().getValue())) {
      return;
    }

    LoginUser loginUser = (LoginUser) context.getPrincipal().getPrincipal();
    claims.claim(SecurityConstants.DETAILS_USER, loginUser);
    claims.subject(loginUser.userInfo().getUserId());
  }
}
