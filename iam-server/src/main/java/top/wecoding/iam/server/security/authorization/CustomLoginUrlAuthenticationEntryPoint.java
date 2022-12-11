package top.wecoding.iam.server.security.authorization;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.web.DefaultRedirectStrategy;
import org.springframework.security.web.RedirectStrategy;
import org.springframework.security.web.authentication.LoginUrlAuthenticationEntryPoint;
import org.springframework.web.util.UriUtils;

import java.nio.charset.StandardCharsets;

/**
 * @author liuyuhui
 * @since 0.5.2
 */
public class CustomLoginUrlAuthenticationEntryPoint extends LoginUrlAuthenticationEntryPoint {

  private final RedirectStrategy redirectStrategy = new DefaultRedirectStrategy();

  /**
   * @param loginFormUrl URL where the login page can be found. Should either be relative to the
   *     web-app context path (include a leading {@code /}) or an absolute URL.
   */
  public CustomLoginUrlAuthenticationEntryPoint(String loginFormUrl) {
    super(loginFormUrl);
  }

  @Override
  protected String determineUrlToUseForThisRequest(
      HttpServletRequest request, HttpServletResponse response, AuthenticationException exception) {
    String redirectUri = request.getParameter("redirect_uri");
    // if (Strings.isBlank(redirectUri)) {
    //   redirectStrategy.sendRedirect(request, response, loginPath);
    //   return;
    // }
    redirectUri = UriUtils.encode(redirectUri, StandardCharsets.UTF_8);
    return String.format("%s?redirect=%s", getLoginFormUrl(), redirectUri);
  }
}
