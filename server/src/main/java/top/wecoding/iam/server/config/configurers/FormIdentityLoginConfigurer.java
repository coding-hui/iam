package top.wecoding.iam.server.config.configurers;

import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import top.wecoding.iam.server.authentication.handler.FormAuthenticationFailureHandler;
import top.wecoding.iam.server.authentication.handler.SsoLogoutSuccessHandler;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
public class FormIdentityLoginConfigurer
    extends AbstractHttpConfigurer<FormIdentityLoginConfigurer, HttpSecurity> {

  @Override
  public void init(HttpSecurity http) throws Exception {
    http.formLogin(
            formLogin -> {
              formLogin.loginPage("/token/login");
              formLogin.loginProcessingUrl("/token/form");
              formLogin.failureHandler(new FormAuthenticationFailureHandler());
            })
        .logout()
        .logoutSuccessHandler(new SsoLogoutSuccessHandler())
        .deleteCookies("JSESSIONID")
        .invalidateHttpSession(true)
        .and()
        .csrf()
        .disable();
  }
}
