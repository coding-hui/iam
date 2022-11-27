package top.wecoding.iam.framework.handler;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.http.HttpHeaders;
import org.springframework.security.core.Authentication;
import org.springframework.security.web.authentication.logout.LogoutSuccessHandler;
import top.wecoding.commons.lang.Strings;

import java.io.IOException;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
public class SsoLogoutSuccessHandler implements LogoutSuccessHandler {

  private static final String REDIRECT_URL = "redirect_url";

  @Override
  public void onLogoutSuccess(
      HttpServletRequest request, HttpServletResponse response, Authentication authentication)
      throws IOException {
    if (response == null) {
      return;
    }

    String redirectUrl = request.getParameter(REDIRECT_URL);
    if (Strings.hasText(redirectUrl)) {
      response.sendRedirect(redirectUrl);
    } else if (Strings.hasText(request.getHeader(HttpHeaders.REFERER))) {
      String referer = request.getHeader(HttpHeaders.REFERER);
      response.sendRedirect(referer);
    }
  }
}
