package top.wecoding.iam.server.security.configurers;

import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import top.wecoding.iam.framework.security.handler.WeCodingAuthenticationFailureEventHandler;
import top.wecoding.iam.server.security.handler.SsoLogoutSuccessHandler;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
public class FormIdentityLoginConfigurer
    extends AbstractHttpConfigurer<FormIdentityLoginConfigurer, HttpSecurity> {

  @Override
  public void init(HttpSecurity http) throws Exception {
    http.formLogin(
        formLogin -> {
          formLogin.loginPage("/auth/login");
          formLogin.loginProcessingUrl("/auth/form");
          formLogin.failureHandler(new WeCodingAuthenticationFailureEventHandler());
        });
    http.logout()
        .logoutSuccessHandler(new SsoLogoutSuccessHandler())
        .deleteCookies("JSESSIONID")
        .invalidateHttpSession(true)
        .and()
        .csrf()
        .disable();
  }
}
