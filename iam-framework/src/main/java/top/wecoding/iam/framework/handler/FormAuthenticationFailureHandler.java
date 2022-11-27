package top.wecoding.iam.framework.handler;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.SneakyThrows;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.web.authentication.AuthenticationFailureHandler;

/**
 * 表单登录失败处理逻辑
 *
 * @author liuyuhui
 * @date 2022/10/3
 */
@Slf4j
public class FormAuthenticationFailureHandler implements AuthenticationFailureHandler {

  @Override
  @SneakyThrows
  public void onAuthenticationFailure(
      HttpServletRequest request, HttpServletResponse response, AuthenticationException exception) {
    log.debug("form authentication failed: {}", exception.getLocalizedMessage());
  }
}
