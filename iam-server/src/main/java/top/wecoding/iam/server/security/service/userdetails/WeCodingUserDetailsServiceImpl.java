package top.wecoding.iam.server.security.service.userdetails;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;
import top.wecoding.commons.core.exception.IllegalParameterException;
import top.wecoding.iam.common.enums.IamErrorCode;
import top.wecoding.iam.common.model.response.UserInfoResponse;
import top.wecoding.iam.common.userdetails.WeCodingUserDetailsService;
import top.wecoding.iam.server.service.UserService;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class WeCodingUserDetailsServiceImpl implements WeCodingUserDetailsService {

  private final UserService userService;

  @Override
  public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
    UserInfoResponse userInfoResponse;
    try {
      userInfoResponse = userService.getInfoByUsername(username);
    } catch (IllegalParameterException e) {
      if (IamErrorCode.USER_DOES_NOT_EXIST.getCode().equalsIgnoreCase(e.getSupplier().getCode())) {
        throw new UsernameNotFoundException(e.getMessage());
      }
      throw e;
    }
    return getUserDetails(userInfoResponse);
  }

  @Override
  public int getOrder() {
    return Integer.MIN_VALUE;
  }
}
