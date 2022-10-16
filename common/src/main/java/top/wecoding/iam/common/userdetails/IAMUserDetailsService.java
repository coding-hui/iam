package top.wecoding.iam.common.userdetails;

import org.springframework.core.Ordered;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.AuthorityUtils;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.util.CollectionUtils;
import top.wecoding.core.result.R;
import top.wecoding.iam.common.constant.SecurityConstants;
import top.wecoding.iam.common.model.GroupInfo;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.common.pojo.UserInfo;

import java.util.Collection;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
public interface IAMUserDetailsService extends UserDetailsService, Ordered {

  default boolean support(String clientId, String grantType) {
    return true;
  }

  /** take one of the biggest */
  default int getOrder() {
    return 0;
  }

  default UserDetails getUserDetails(R<UserInfoResponse> result) {
    UserInfoResponse userInfoResponse = result.getData();
    UserInfo info = userInfoResponse.getUserInfo();
    List<GroupInfo> groups = userInfoResponse.getGroups();
    Set<String> roles = userInfoResponse.getRoles();
    Set<String> permissions = userInfoResponse.getPermissions();

    Set<String> authsSet = new HashSet<>();
    if (!CollectionUtils.isEmpty(permissions)) {
      authsSet.addAll(permissions);
    }
    if (!CollectionUtils.isEmpty(roles)) {
      authsSet.addAll(
          roles.stream()
              .map(role -> SecurityConstants.ROLE_PREFIX + role)
              .collect(Collectors.toList()));
    }
    Collection<GrantedAuthority> authorities =
        AuthorityUtils.createAuthorityList(authsSet.toArray(new String[0]));

    return new LoginUser(authorities, info, groups, permissions, roles);
  }

  default UserDetails loadUserByUser(LoginUser loginUser) {
    return this.loadUserByUsername(loginUser.getUsername());
  }
}
