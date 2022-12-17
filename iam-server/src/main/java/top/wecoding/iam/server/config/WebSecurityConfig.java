package top.wecoding.iam.server.config;

import jakarta.annotation.Resource;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.web.SecurityFilterChain;
import top.wecoding.iam.server.security.authorization.authentication.WeCodingDaoAuthenticationProvider;
import top.wecoding.iam.server.security.configurers.FormIdentityLoginConfigurer;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
@EnableWebSecurity
@RequiredArgsConstructor
@Configuration(proxyBeanMethods = false)
public class WebSecurityConfig {

  @Resource private final WeCodingDaoAuthenticationProvider weCodingDaoAuthenticationProvider;

  @Bean
  SecurityFilterChain defaultSecurityFilterChain(HttpSecurity http) throws Exception {
    http.authorizeHttpRequests(
            authorize ->
                authorize.requestMatchers("/login").permitAll().anyRequest().authenticated())
        .headers()
        .frameOptions()
        .sameOrigin()
        .and()
        .csrf()
        .disable();

    http.apply(new FormIdentityLoginConfigurer());
    http.authenticationProvider(weCodingDaoAuthenticationProvider);
    return http.build();
  }
}
