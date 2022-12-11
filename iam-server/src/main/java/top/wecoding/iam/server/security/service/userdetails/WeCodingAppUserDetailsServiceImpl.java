package top.wecoding.iam.server.security.service.userdetails;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.cache.CacheManager;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import top.wecoding.iam.common.constant.SecurityConstants;
import top.wecoding.iam.common.userdetails.WeCodingUserDetailsService;
import top.wecoding.iam.server.service.UserService;

/**
 * @author liuyuhui
 * @date 2022/10/4
 */
@Slf4j
@RequiredArgsConstructor
public class WeCodingAppUserDetailsServiceImpl implements WeCodingUserDetailsService {

  private final UserService userService;

  private final CacheManager cacheManager;

  @Override
  public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
    return null;
  }

  @Override
  public boolean support(String clientId, String grantType) {
    return SecurityConstants.APP.equals(clientId);
  }
}
