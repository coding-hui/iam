package top.wecoding.iam.server.config;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Bean;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.oauth2.server.resource.introspection.OpaqueTokenIntrospector;
import org.springframework.security.oauth2.server.resource.web.BearerTokenResolver;
import org.springframework.security.web.SecurityFilterChain;
import top.wecoding.iam.framework.props.IgnoreWhiteProperties;
import top.wecoding.iam.framework.security.web.ResourceAuthExceptionEntryPoint;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
@Slf4j
@RequiredArgsConstructor
public class ResourceServerConfig {

  protected final ResourceAuthExceptionEntryPoint resourceAuthExceptionEntryPoint;

  private final IgnoreWhiteProperties permitAllUrl;

  private final BearerTokenResolver weCodingBearerTokenExtractor;

  private final OpaqueTokenIntrospector opaqueTokenIntrospector;

  @Bean
  @Order(Ordered.HIGHEST_PRECEDENCE + 1)
  SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {

    http.authorizeHttpRequests(
            authorize ->
                authorize
                    .requestMatchers(permitAllUrl.getWhites().toArray(String[]::new))
                    .permitAll()
                    .anyRequest()
                    .authenticated())
        .oauth2ResourceServer(
            oauth2 ->
                oauth2
                    .opaqueToken(token -> token.introspector(opaqueTokenIntrospector))
                    .authenticationEntryPoint(resourceAuthExceptionEntryPoint)
                    .bearerTokenResolver(weCodingBearerTokenExtractor))
        .headers()
        .frameOptions()
        .disable()
        .and()
        .csrf()
        .disable();

    return http.build();
  }
}
