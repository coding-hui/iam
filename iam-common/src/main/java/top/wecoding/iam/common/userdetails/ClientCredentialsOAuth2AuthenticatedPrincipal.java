package top.wecoding.iam.common.userdetails;

import lombok.RequiredArgsConstructor;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.oauth2.core.OAuth2AuthenticatedPrincipal;

import java.util.Collection;
import java.util.Map;

/**
 * @author liuyuhui
 * @since 0.5
 */
@RequiredArgsConstructor
public class ClientCredentialsOAuth2AuthenticatedPrincipal implements OAuth2AuthenticatedPrincipal {

  private final Map<String, Object> attributes;

  private final Collection<GrantedAuthority> authorities;

  private final String name;

  @Override
  public Map<String, Object> getAttributes() {
    return this.attributes;
  }

  @Override
  public Collection<? extends GrantedAuthority> getAuthorities() {
    return this.authorities;
  }

  @Override
  public String getName() {
    return this.name;
  }
}
