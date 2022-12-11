package top.wecoding.iam.server.config;

import jakarta.annotation.Resource;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.web.SecurityFilterChain;
import top.wecoding.iam.server.security.authorization.authentication.WeCodingDaoAuthenticationProvider;
import top.wecoding.iam.server.security.configurers.FormIdentityLoginConfigurer;

import static org.springframework.security.config.Customizer.withDefaults;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
@EnableWebSecurity
@RequiredArgsConstructor
@Configuration(proxyBeanMethods = false)
public class WebSecurityConfig {

  @Resource private final UserDetailsService userDetailsService;

  @Resource private final WeCodingDaoAuthenticationProvider weCodingDaoAuthenticationProvider;

  @Bean
  @Order(Ordered.HIGHEST_PRECEDENCE + 2)
  SecurityFilterChain defaultSecurityFilterChain(HttpSecurity http) throws Exception {
    http.authorizeHttpRequests(authorize -> authorize.anyRequest().authenticated())
        .formLogin(withDefaults())
        .headers()
        .frameOptions()
        .disable()
        .and()
        .csrf()
        .disable();

    http.apply(new FormIdentityLoginConfigurer());

    http.authenticationProvider(weCodingDaoAuthenticationProvider)
        .userDetailsService(userDetailsService);
    return http.build();
  }
}
