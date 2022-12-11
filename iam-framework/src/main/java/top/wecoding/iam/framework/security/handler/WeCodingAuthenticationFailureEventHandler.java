package top.wecoding.iam.framework.security.handler;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.SneakyThrows;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.converter.json.MappingJackson2HttpMessageConverter;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.oauth2.core.endpoint.OAuth2ParameterNames;
import org.springframework.security.web.authentication.AuthenticationFailureHandler;
import top.wecoding.commons.core.util.ServletResponseUtil;
import top.wecoding.iam.common.enums.IamErrorCode;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
@Slf4j
public class WeCodingAuthenticationFailureEventHandler implements AuthenticationFailureHandler {

  private final MappingJackson2HttpMessageConverter errorHttpResponseConverter =
      new MappingJackson2HttpMessageConverter();

  /**
   * Called when an authentication attempt fails.
   *
   * @param request the request during which the authentication attempt occurred.
   * @param response the response.
   * @param exception the exception which was thrown to reject the authentication request.
   */
  @Override
  @SneakyThrows
  public void onAuthenticationFailure(
      HttpServletRequest request, HttpServletResponse response, AuthenticationException exception) {
    String username = request.getParameter(OAuth2ParameterNames.USERNAME);

    log.info("user: {} authentication failed：{}", username, exception.getLocalizedMessage());

    ServletResponseUtil.webMvcResponseWriter(
        response,
        HttpStatus.UNAUTHORIZED,
        IamErrorCode.UNAUTHORIZED,
        exception.getLocalizedMessage());
  }
}
