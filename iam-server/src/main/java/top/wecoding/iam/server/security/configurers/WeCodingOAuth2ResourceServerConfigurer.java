package top.wecoding.iam.server.security.configurers;

import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configurers.oauth2.server.resource.OAuth2ResourceServerConfigurer;
import org.springframework.security.oauth2.server.resource.introspection.OpaqueTokenIntrospector;
import org.springframework.security.oauth2.server.resource.web.BearerTokenResolver;
import top.wecoding.iam.framework.props.IgnoreWhiteProperties;
import top.wecoding.iam.framework.security.web.ResourceAuthExceptionEntryPoint;

/**
 * @author liuyuhui
 * @since 0.5
 */
@RequiredArgsConstructor
@Configuration(proxyBeanMethods = false)
public class WeCodingOAuth2ResourceServerConfigurer {

  private final IgnoreWhiteProperties permitAllUrl;

  private final OpaqueTokenIntrospector opaqueTokenIntrospector;

  private final BearerTokenResolver weCodingBearerTokenExtractor;

  private final ResourceAuthExceptionEntryPoint resourceAuthExceptionEntryPoint;

  @Bean
  Customizer<OAuth2ResourceServerConfigurer<HttpSecurity>> oauth2ResourceServerCustomizer() {
    return customizer ->
        customizer
            .opaqueToken(token -> token.introspector(opaqueTokenIntrospector))
            .authenticationEntryPoint(resourceAuthExceptionEntryPoint)
            .bearerTokenResolver(weCodingBearerTokenExtractor);
  }
}
