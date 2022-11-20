package top.wecoding.iam.server.config;

import com.nimbusds.jose.jwk.JWKSet;
import com.nimbusds.jose.jwk.RSAKey;
import com.nimbusds.jose.jwk.source.JWKSource;
import com.nimbusds.jose.proc.SecurityContext;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.OAuth2AuthorizationServerConfiguration;
import org.springframework.security.config.annotation.web.configurers.oauth2.server.authorization.OAuth2AuthorizationServerConfigurer;
import org.springframework.security.config.annotation.web.configurers.oauth2.server.resource.OAuth2ResourceServerConfigurer;
import org.springframework.security.oauth2.jwt.JwtDecoder;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.security.oauth2.server.authorization.config.ProviderSettings;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.util.matcher.RequestMatcher;
import top.wecoding.iam.common.constant.SecurityConstants;
import top.wecoding.iam.server.jose.Jwks;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@SuppressWarnings("all")
@RequiredArgsConstructor
@Configuration(proxyBeanMethods = false)
public class AuthorizationServerConfig {

  private final OAuth2AuthorizationService authorizationService;

  @Bean
  @Order(Ordered.HIGHEST_PRECEDENCE)
  public SecurityFilterChain authorizationServerSecurityFilterChain(HttpSecurity http)
      throws Exception {
    OAuth2AuthorizationServerConfigurer<HttpSecurity> authorizationServerConfigurer =
        new OAuth2AuthorizationServerConfigurer<>();

    RequestMatcher endpointsMatcher = authorizationServerConfigurer.getEndpointsMatcher();

    http.requestMatcher(endpointsMatcher)
        .authorizeHttpRequests(authorize -> authorize.anyRequest().authenticated())
        .csrf(csrf -> csrf.ignoringRequestMatchers(endpointsMatcher))
        .apply(new FormIdentityLoginConfigurer())
        .and()
        .apply(authorizationServerConfigurer)
        .and()
        .oauth2ResourceServer(OAuth2ResourceServerConfigurer::jwt);
    return http.build();
  }

  @Bean
  public ProviderSettings providerSettings() {
    return ProviderSettings.builder().issuer(SecurityConstants.PROJECT_LICENSE).build();
  }

  @Bean
  public JWKSource<SecurityContext> jwkSource() {
    RSAKey rsaKey = Jwks.generateRsa();
    JWKSet jwkSet = new JWKSet(rsaKey);
    return (jwkSelector, securityContext) -> jwkSelector.select(jwkSet);
  }

  @Bean
  public JwtDecoder jwtDecoder(JWKSource<SecurityContext> jwkSource) {
    return OAuth2AuthorizationServerConfiguration.jwtDecoder(jwkSource);
  }

  // public OAuth2TokenGenerator oAuth2TokenGenerator() {
  //   IAMOAuth2AccessTokenGenerator accessTokenGenerator = new IAMOAuth2AccessTokenGenerator();
  //   accessTokenGenerator.setAccessTokenCustomizer(new IAMOAuth2TokenCustomizer());
  //   return new DelegatingOAuth2TokenGenerator(
  //       accessTokenGenerator, new IAMOAuth2RefreshTokenGenerator());
  // }

  // /**
  //  * request -> xToken 注入请求转换器
  //  *
  //  * @return DelegatingAuthenticationConverter
  //  */
  // private AuthenticationConverter accessTokenRequestConverter() {
  //   return new DelegatingAuthenticationConverter(
  //       Arrays.asList(
  //           new OAuth2ResourceOwnerPasswordAuthenticationConverter(),
  //           // new OAuth2ResourceOwnerSmsAuthenticationConverter(),
  //           new OAuth2RefreshTokenAuthenticationConverter(),
  //           new OAuth2ClientCredentialsAuthenticationConverter(),
  //           new OAuth2AuthorizationCodeAuthenticationConverter(),
  //           new OAuth2AuthorizationCodeRequestAuthenticationConverter()));
  // }

  // @SuppressWarnings("unchecked")
  // private void addCustomOAuth2GrantAuthenticationProvider(HttpSecurity http) {
  //   AuthenticationManager authenticationManager =
  // http.getSharedObject(AuthenticationManager.class);
  //   OAuth2AuthorizationService authorizationService =
  //       http.getSharedObject(OAuth2AuthorizationService.class);
  //
  //   OAuth2ResourceOwnerPasswordAuthenticationProvider resourceOwnerPasswordAuthenticationProvider
  // =
  //       new OAuth2ResourceOwnerPasswordAuthenticationProvider(
  //           authenticationManager, authorizationService, oAuth2TokenGenerator());
  //
  //   // 处理 UsernamePasswordAuthenticationToken
  //   http.authenticationProvider(new IAMDaoAuthenticationProvider());
  //   // 处理 OAuth2ResourceOwnerPasswordAuthenticationToken
  //   http.authenticationProvider(resourceOwnerPasswordAuthenticationProvider);
  // }
}
