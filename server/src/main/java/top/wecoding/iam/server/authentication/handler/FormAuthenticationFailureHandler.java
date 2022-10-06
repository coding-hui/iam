package top.wecoding.iam.server.authentication.handler;

import cn.hutool.core.util.CharsetUtil;
import cn.hutool.http.HttpUtil;
import lombok.SneakyThrows;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.web.authentication.AuthenticationFailureHandler;
import top.wecoding.core.util.WebUtils;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

/**
 * 表单登录失败处理逻辑
 *
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Slf4j
public class FormAuthenticationFailureHandler implements AuthenticationFailureHandler {

  @Override
  @SneakyThrows
  public void onAuthenticationFailure(
      HttpServletRequest request, HttpServletResponse response, AuthenticationException exception) {
    log.debug("form authentication failed: {}", exception.getLocalizedMessage());
    String url =
        HttpUtil.encodeParams(
            String.format("/token/login?error=%s", exception.getMessage()),
            CharsetUtil.CHARSET_UTF_8);
    WebUtils.getResponse().sendRedirect(url);
  }
}
