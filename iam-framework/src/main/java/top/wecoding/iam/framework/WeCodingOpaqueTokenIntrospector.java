package top.wecoding.iam.framework;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.authority.AuthorityUtils;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.oauth2.core.AuthorizationGrantType;
import org.springframework.security.oauth2.core.OAuth2AuthenticatedPrincipal;
import org.springframework.security.oauth2.server.authorization.OAuth2Authorization;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.security.oauth2.server.authorization.OAuth2TokenType;
import org.springframework.security.oauth2.server.resource.InvalidBearerTokenException;
import org.springframework.security.oauth2.server.resource.introspection.OpaqueTokenIntrospector;
import org.springframework.stereotype.Component;
import top.wecoding.iam.common.userdetails.ClientCredentialsOAuth2AuthenticatedPrincipal;
import top.wecoding.iam.common.userdetails.LoginUser;

import java.security.Principal;
import java.util.Objects;

/**
 * @author liuyuhui
 * @since 0.5
 */
@Slf4j
@Component
@RequiredArgsConstructor
public class WeCodingOpaqueTokenIntrospector implements OpaqueTokenIntrospector {

  private final OAuth2AuthorizationService authorizationService;

  private final UserDetailsService userDetailsService;

  @Override
  public OAuth2AuthenticatedPrincipal introspect(String token) {
    OAuth2Authorization oldAuthorization =
        authorizationService.findByToken(token, OAuth2TokenType.ACCESS_TOKEN);
    if (Objects.isNull(oldAuthorization)) {
      throw new InvalidBearerTokenException(token);
    }

    // 客户端模式默认返回
    if (AuthorizationGrantType.CLIENT_CREDENTIALS.equals(
        oldAuthorization.getAuthorizationGrantType())) {
      return new ClientCredentialsOAuth2AuthenticatedPrincipal(
          oldAuthorization.getAttributes(),
          AuthorityUtils.NO_AUTHORITIES,
          oldAuthorization.getPrincipalName());
    }

    UserDetails userDetails = null;
    try {
      Object principal =
          Objects.requireNonNull(oldAuthorization).getAttributes().get(Principal.class.getName());
      UsernamePasswordAuthenticationToken usernamePasswordAuthenticationToken =
          (UsernamePasswordAuthenticationToken) principal;
      Object tokenPrincipal = usernamePasswordAuthenticationToken.getPrincipal();
      userDetails =
          userDetailsService.loadUserByUsername(((LoginUser) tokenPrincipal).getUsername());
    } catch (UsernameNotFoundException notFoundException) {
      log.warn("用户不不存在 {}", notFoundException.getLocalizedMessage());
      throw notFoundException;
    } catch (Exception ex) {
      log.error("资源服务器 introspect Token error {}", ex.getLocalizedMessage());
    }
    return (LoginUser) userDetails;
  }
}
