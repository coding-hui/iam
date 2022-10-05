package top.wecoding.iam.sdk.introspection;

import cn.hutool.extra.spring.SpringUtil;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.core.Ordered;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.authority.AuthorityUtils;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.oauth2.core.AuthorizationGrantType;
import org.springframework.security.oauth2.core.OAuth2AuthenticatedPrincipal;
import org.springframework.security.oauth2.core.OAuth2TokenType;
import org.springframework.security.oauth2.server.authorization.OAuth2Authorization;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.security.oauth2.server.resource.InvalidBearerTokenException;
import org.springframework.security.oauth2.server.resource.introspection.OpaqueTokenIntrospector;
import top.wecoding.iam.common.userdetails.IAMUserDetailsService;
import top.wecoding.iam.common.userdetails.LoginUser;
import top.wecoding.iam.common.userdetails.IAMClientCredentialsOAuth2AuthenticatedPrincipal;

import java.security.Principal;
import java.util.Comparator;
import java.util.Map;
import java.util.Objects;
import java.util.Optional;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Slf4j
@RequiredArgsConstructor
public class IAMOpaqueTokenIntrospector implements OpaqueTokenIntrospector {

  private final OAuth2AuthorizationService authorizationService;

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
      return new IAMClientCredentialsOAuth2AuthenticatedPrincipal(
          oldAuthorization.getAttributes(),
          AuthorityUtils.NO_AUTHORITIES,
          oldAuthorization.getPrincipalName());
    }

    Map<String, IAMUserDetailsService> userDetailsServiceMap =
        SpringUtil.getBeansOfType(IAMUserDetailsService.class);

    Optional<IAMUserDetailsService> optional =
        userDetailsServiceMap.values().stream()
            .filter(
                service ->
                    service.support(
                        Objects.requireNonNull(oldAuthorization).getRegisteredClientId(),
                        oldAuthorization.getAuthorizationGrantType().getValue()))
            .max(Comparator.comparingInt(Ordered::getOrder));

    UserDetails userDetails = null;
    try {
      Object principal =
          Objects.requireNonNull(oldAuthorization).getAttributes().get(Principal.class.getName());
      UsernamePasswordAuthenticationToken usernamePasswordAuthenticationToken =
          (UsernamePasswordAuthenticationToken) principal;
      Object tokenPrincipal = usernamePasswordAuthenticationToken.getPrincipal();
      userDetails = optional.get().loadUserByUser((LoginUser) tokenPrincipal);
    } catch (UsernameNotFoundException notFoundException) {
      log.warn("用户不不存在 {}", notFoundException.getLocalizedMessage());
      throw notFoundException;
    } catch (Exception ex) {
      log.error("资源服务器 introspect Token error {}", ex.getLocalizedMessage());
    }
    return (LoginUser) userDetails;
  }
}
