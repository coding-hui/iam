package top.wecoding.iam.server.oauth2.userdetails;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.cache.Cache;
import org.springframework.cache.CacheManager;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import top.wecoding.iam.common.constant.RedisConstant;
import top.wecoding.iam.common.constant.SecurityConstants;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.common.userdetails.IAMUserDetailsService;
import top.wecoding.iam.common.userdetails.LoginUser;
import top.wecoding.iam.server.service.UserService;

/**
 * @author liuyuhui
 * @date 2022/10/4
 * @qq 1515418211
 */
@Slf4j
@RequiredArgsConstructor
public class IAMAppUserDetailsServiceImpl implements IAMUserDetailsService {

  private final UserService userService;

  private final CacheManager cacheManager;

  @Override
  public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
    Cache cache = cacheManager.getCache(RedisConstant.USER_DETAILS);
    if (cache != null && cache.get(username) != null) {
      return (LoginUser) cache.get(username).get();
    }

    UserInfoResponse userInfoResponse = userService.getInfoByUsername(username);
    UserDetails userDetails = getUserDetails(userInfoResponse);
    if (cache != null) {
      cache.put(username, userDetails);
    }
    return userDetails;
  }

  @Override
  public boolean support(String clientId, String grantType) {
    return SecurityConstants.APP.equals(clientId);
  }
}
