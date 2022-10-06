package top.wecoding.iam.common.userdetails;

import org.springframework.core.Ordered;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.AuthorityUtils;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.core.result.R;
import top.wecoding.iam.common.model.UserInfo;
import top.wecoding.iam.common.model.response.UserInfoResponse;

import java.util.Collection;
import java.util.HashSet;
import java.util.Set;

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

    Set<String> dbAuthsSet = new HashSet<>();
    Collection<GrantedAuthority> authorities =
        AuthorityUtils.createAuthorityList(dbAuthsSet.toArray(new String[0]));

    return new LoginUser(
        info.getUsername(),
        SecurityConstants.BCRYPT + info.getPassword(),
        true,
        true,
        true,
        true,
        authorities,
        info.getUserId(),
        info.getTenantId(),
        info.getPhone());
  }

  default UserDetails loadUserByUser(LoginUser loginUser) {
    return this.loadUserByUsername(loginUser.getUsername());
  }
}
