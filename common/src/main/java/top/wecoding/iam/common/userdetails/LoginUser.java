package top.wecoding.iam.common.userdetails;

import lombok.Getter;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.SpringSecurityCoreVersion;
import org.springframework.security.core.userdetails.User;
import org.springframework.security.oauth2.core.OAuth2AuthenticatedPrincipal;

import java.util.Collection;
import java.util.HashMap;
import java.util.Map;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
public class LoginUser extends User implements OAuth2AuthenticatedPrincipal {

  private static final long serialVersionUID = SpringSecurityCoreVersion.SERIAL_VERSION_UID;

  @Getter private final String userId;

  @Getter private final String tenantId;

  @Getter private final String phone;

  public LoginUser(
      String username,
      String password,
      boolean enabled,
      boolean accountNonExpired,
      boolean credentialsNonExpired,
      boolean accountNonLocked,
      Collection<? extends GrantedAuthority> authorities,
      String userId,
      String tenantId,
      String phone) {
    super(
        username,
        password,
        enabled,
        accountNonExpired,
        credentialsNonExpired,
        accountNonLocked,
        authorities);
    this.userId = userId;
    this.tenantId = tenantId;
    this.phone = phone;
  }

  @Override
  public Map<String, Object> getAttributes() {
    return new HashMap<>();
  }

  @Override
  public String getName() {
    return this.getUsername();
  }
}
