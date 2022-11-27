package top.wecoding.iam.server.config;

import jakarta.annotation.Resource;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.oauth2.server.resource.OAuth2ResourceServerConfigurer;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.web.SecurityFilterChain;
import top.wecoding.iam.framework.authentication.WeCodingDaoAuthenticationProvider;

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
  SecurityFilterChain defaultSecurityFilterChain(HttpSecurity http) throws Exception {
    http.authorizeHttpRequests(authorize -> authorize.anyRequest().authenticated())
        .formLogin(withDefaults())
        .oauth2ResourceServer(OAuth2ResourceServerConfigurer::jwt)
        .headers()
        .frameOptions()
        .disable()
        .and()
        .csrf()
        .disable();

    http.authenticationProvider(weCodingDaoAuthenticationProvider)
        .userDetailsService(userDetailsService);
    return http.build();
  }
}
