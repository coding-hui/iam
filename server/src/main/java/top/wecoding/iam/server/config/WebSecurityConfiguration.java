package top.wecoding.iam.server.config;

import org.springframework.context.annotation.Bean;
import org.springframework.core.annotation.Order;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.web.SecurityFilterChain;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@EnableWebSecurity(debug = true)
public class WebSecurityConfiguration {

  @Bean
  public SecurityFilterChain defaultSecurityFilterChain(HttpSecurity http) throws Exception {
    http.authorizeRequests(
            authorizeRequests ->
                authorizeRequests.antMatchers("/token/*").permitAll().anyRequest().authenticated())
        .headers()
        .frameOptions()
        // 避免iframe同源无法登录
        .sameOrigin()
        .and()
        .apply(new FormIdentityLoginConfigurer());
    return http.build();
  }

  @Bean
  @Order(0)
  public SecurityFilterChain resources(HttpSecurity http) throws Exception {
    http.requestMatchers((matchers) -> matchers.antMatchers("/actuator/**", "/css/**", "/error"))
        .authorizeHttpRequests((authorize) -> authorize.anyRequest().permitAll())
        .requestCache()
        .disable()
        .securityContext()
        .disable()
        .sessionManagement()
        .disable();
    return http.build();
  }
}
