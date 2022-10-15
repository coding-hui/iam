package top.wecoding.iam.sdk.userdetails;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Primary;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.core.result.R;
import top.wecoding.iam.api.feign.RemoteUserService;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.common.userdetails.IAMUserDetailsService;
import top.wecoding.iam.sdk.cache.UserDetailsCacheKeyBuilder;
import top.wecoding.redis.util.RedisUtils;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Slf4j
@Primary
@RequiredArgsConstructor
public class IAMUserDetailsServiceImpl implements IAMUserDetailsService {

  private final RemoteUserService remoteUserService;

  @Override
  public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
    return RedisUtils.get(
        new UserDetailsCacheKeyBuilder().key(username),
        () -> {
          R<UserInfoResponse> result = remoteUserService.info(username, SecurityConstants.INNER);
          return getUserDetails(result);
        });
  }

  @Override
  public int getOrder() {
    return Integer.MIN_VALUE;
  }
}
