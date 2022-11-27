package top.wecoding.iam.framework.configurer;

import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import top.wecoding.iam.framework.handler.SsoLogoutSuccessHandler;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
public class FormIdentityLoginConfigurer
    extends AbstractHttpConfigurer<FormIdentityLoginConfigurer, HttpSecurity> {

  @Override
  public void init(HttpSecurity http) throws Exception {
    http.formLogin(Customizer.withDefaults())
        .logout()
        .logoutSuccessHandler(new SsoLogoutSuccessHandler())
        .deleteCookies("JSESSIONID")
        .invalidateHttpSession(true)
        .and()
        .csrf()
        .disable();
  }
}
