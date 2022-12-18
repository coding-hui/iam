package top.wecoding.iam.server.config;

import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.oauth2.server.resource.OAuth2ResourceServerConfigurer;
import org.springframework.security.web.SecurityFilterChain;
import top.wecoding.iam.framework.props.IgnoreWhiteProperties;
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

  private final IgnoreWhiteProperties permitAllUrl;

  private final WeCodingDaoAuthenticationProvider authenticationProvider;

  @Bean
  SecurityFilterChain defaultSecurityFilterChain(
      HttpSecurity http,
      Customizer<OAuth2ResourceServerConfigurer<HttpSecurity>> oauth2ResourceServerCustomizer)
      throws Exception {
    http.authorizeHttpRequests(
            authorize ->
                authorize
                    .requestMatchers(permitAllUrl.getWhites().toArray(new String[] {}))
                    .permitAll()
                    .anyRequest()
                    .authenticated())
        .headers()
        .frameOptions()
        .sameOrigin()
        .and()
        .csrf()
        .disable();

    http.apply(new FormIdentityLoginConfigurer())
        .and()
        .authenticationProvider(authenticationProvider);

    http.oauth2ResourceServer(oauth2ResourceServerCustomizer);

    return http.build();
  }
}
