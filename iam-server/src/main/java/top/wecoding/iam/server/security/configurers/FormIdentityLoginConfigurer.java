package top.wecoding.iam.server.security.configurers;

import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import top.wecoding.iam.server.security.handler.SsoLogoutSuccessHandler;

import static org.springframework.security.config.Customizer.withDefaults;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
public class FormIdentityLoginConfigurer
    extends AbstractHttpConfigurer<FormIdentityLoginConfigurer, HttpSecurity> {

  @Override
  public void init(HttpSecurity http) throws Exception {
    http.formLogin(withDefaults());
    http.logout()
        .logoutSuccessHandler(new SsoLogoutSuccessHandler())
        .deleteCookies("JSESSIONID")
        .invalidateHttpSession(true)
        .and()
        .csrf()
        .disable();
  }
}
