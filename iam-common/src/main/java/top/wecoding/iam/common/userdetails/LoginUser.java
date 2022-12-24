package top.wecoding.iam.common.userdetails;

import java.io.Serial;
import java.util.*;
import lombok.Getter;
import lombok.experimental.Accessors;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.SpringSecurityCoreVersion;
import org.springframework.security.core.userdetails.User;
import org.springframework.security.oauth2.core.OAuth2AuthenticatedPrincipal;
import top.wecoding.iam.common.constant.SecurityConstants;
import top.wecoding.iam.common.entity.UserInfo;
import top.wecoding.iam.common.model.GroupInfo;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
@Getter
@Accessors(fluent = true)
public class LoginUser extends User implements OAuth2AuthenticatedPrincipal {

  @Serial private static final long serialVersionUID = SpringSecurityCoreVersion.SERIAL_VERSION_UID;

  private final UserInfo userInfo;

  private final List<GroupInfo> groups;

  private final Set<String> permissions;

  private final Set<String> roles;

  public LoginUser(
      Collection<? extends GrantedAuthority> authorities,
      UserInfo userInfo,
      List<GroupInfo> groups,
      Set<String> permissions,
      Set<String> roles) {
    super(
        userInfo.getUsername(),
        SecurityConstants.BCRYPT + userInfo.getPassword(),
        true,
        true,
        true,
        true,
        authorities);
    this.userInfo = userInfo;
    this.groups = groups;
    this.permissions = permissions;
    this.roles = roles;
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
