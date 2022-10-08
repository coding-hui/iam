package top.wecoding.iam.common.userdetails;

import java.util.Collection;
import java.util.Map;
import lombok.RequiredArgsConstructor;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.oauth2.core.OAuth2AuthenticatedPrincipal;

/**
 * 支持客户端模式的用户存储
 *
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@RequiredArgsConstructor
public class IAMClientCredentialsOAuth2AuthenticatedPrincipal
    implements OAuth2AuthenticatedPrincipal {

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
