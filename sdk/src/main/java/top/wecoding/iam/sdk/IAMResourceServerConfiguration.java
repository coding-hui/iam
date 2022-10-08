package top.wecoding.iam.sdk;

import cn.hutool.core.util.ArrayUtil;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Bean;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.oauth2.server.resource.introspection.OpaqueTokenIntrospector;
import org.springframework.security.web.SecurityFilterChain;
import top.wecoding.iam.sdk.web.IAMBearerTokenExtractor;
import top.wecoding.iam.sdk.web.ResourceAuthExceptionEntryPoint;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Slf4j
@EnableWebSecurity
@RequiredArgsConstructor
public class IAMResourceServerConfiguration {

  protected final ResourceAuthExceptionEntryPoint resourceAuthExceptionEntryPoint;

  private final IgnoreWhiteProperties permitAllUrl;

  private final IAMBearerTokenExtractor iamBearerTokenExtractor;

  private final OpaqueTokenIntrospector customOpaqueTokenIntrospector;

  @Bean
  @Order(Ordered.HIGHEST_PRECEDENCE)
  SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {

    http.authorizeRequests(
            authorizeRequests ->
                authorizeRequests
                    .antMatchers(ArrayUtil.toArray(permitAllUrl.getWhites(), String.class))
                    .permitAll()
                    .anyRequest()
                    .authenticated())
        .oauth2ResourceServer(
            oauth2 ->
                oauth2
                    .opaqueToken(token -> token.introspector(customOpaqueTokenIntrospector))
                    .authenticationEntryPoint(resourceAuthExceptionEntryPoint)
                    .bearerTokenResolver(iamBearerTokenExtractor))
        .headers()
        .frameOptions()
        .disable()
        .and()
        .csrf()
        .disable();

    return http.build();
  }
}
