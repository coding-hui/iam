package top.wecoding.iam.server.config;

import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
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
    http.formLogin(Customizer.withDefaults())
        .logout()
        .logoutSuccessHandler(new SsoLogoutSuccessHandler())
        .and()
        .csrf()
        .disable();
  }
}
