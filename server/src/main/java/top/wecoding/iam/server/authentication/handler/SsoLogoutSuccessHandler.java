package top.wecoding.iam.server.authentication.handler;

import cn.hutool.core.util.StrUtil;
import java.io.IOException;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import org.springframework.http.HttpHeaders;
import org.springframework.security.core.Authentication;
import org.springframework.security.web.authentication.logout.LogoutSuccessHandler;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
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
    if (StrUtil.isNotBlank(redirectUrl)) {
      response.sendRedirect(redirectUrl);
    } else if (StrUtil.isNotBlank(request.getHeader(HttpHeaders.REFERER))) {
      String referer = request.getHeader(HttpHeaders.REFERER);
      response.sendRedirect(referer);
    }
  }
}
