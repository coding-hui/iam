package top.wecoding.iam.server.oauth2.userdetails;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.common.userdetails.IAMUserDetailsService;
import top.wecoding.iam.framework.cache.UserDetailsCacheKeyBuilder;
import top.wecoding.iam.server.service.UserService;
import top.wecoding.redis.util.RedisUtils;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class IAMUserDetailsServiceImpl implements IAMUserDetailsService {

  private final UserService userService;

  @Override
  public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
    return RedisUtils.get(
        new UserDetailsCacheKeyBuilder().key(username),
        () -> {
          UserInfoResponse userInfoResponse = userService.getInfoByUsername(username);
          return getUserDetails(userInfoResponse);
        });
  }

  @Override
  public int getOrder() {
    return Integer.MIN_VALUE;
  }
}
