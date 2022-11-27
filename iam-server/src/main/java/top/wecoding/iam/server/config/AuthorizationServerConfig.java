package top.wecoding.iam.server.config;

import com.nimbusds.jose.jwk.JWKSet;
import com.nimbusds.jose.jwk.RSAKey;
import com.nimbusds.jose.jwk.source.JWKSource;
import com.nimbusds.jose.proc.SecurityContext;
import jakarta.annotation.Resource;
import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.oauth2.core.OAuth2Token;
import org.springframework.security.oauth2.jwt.JwtDecoder;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.security.oauth2.server.authorization.config.annotation.web.configuration.OAuth2AuthorizationServerConfiguration;
import org.springframework.security.oauth2.server.authorization.config.annotation.web.configurers.OAuth2AuthorizationServerConfigurer;
import org.springframework.security.oauth2.server.authorization.settings.AuthorizationServerSettings;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenGenerator;
import org.springframework.security.oauth2.server.authorization.web.authentication.*;
import org.springframework.security.web.DefaultSecurityFilterChain;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.AuthenticationConverter;
import org.springframework.security.web.util.matcher.RequestMatcher;
import top.wecoding.iam.framework.authentication.WeCodingDaoAuthenticationProvider;
import top.wecoding.iam.framework.authentication.password.OAuth2ResourceOwnerPasswordAuthenticationConverter;
import top.wecoding.iam.framework.authentication.password.OAuth2ResourceOwnerPasswordAuthenticationProvider;
import top.wecoding.iam.framework.configurer.FormIdentityLoginConfigurer;
import top.wecoding.iam.framework.handler.WeCodingAuthenticationFailureEventHandler;
import top.wecoding.iam.framework.jose.Jwks;

import java.util.Arrays;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
@RequiredArgsConstructor
@Configuration(proxyBeanMethods = false)
public class AuthorizationServerConfig {

  private final OAuth2AuthorizationService authorizationService;

  @Resource private final WeCodingDaoAuthenticationProvider weCodingDaoAuthenticationProvider;

  public static void applyDefaultSecurity(HttpSecurity http) throws Exception {
    OAuth2AuthorizationServerConfigurer authorizationServerConfigurer =
        new OAuth2AuthorizationServerConfigurer();
    RequestMatcher endpointsMatcher = authorizationServerConfigurer.getEndpointsMatcher();

    http.securityMatcher(endpointsMatcher)
        .authorizeHttpRequests(authorize -> authorize.anyRequest().authenticated())
        .csrf(csrf -> csrf.ignoringRequestMatchers(endpointsMatcher))
        .apply(authorizationServerConfigurer);
  }

  @Bean
  @Order(Ordered.HIGHEST_PRECEDENCE)
  public SecurityFilterChain authorizationServerSecurityFilterChain(HttpSecurity http)
      throws Exception {
    applyDefaultSecurity(http);

    http.getConfigurer(OAuth2AuthorizationServerConfigurer.class)
        .tokenEndpoint(
            tokenEndpoint -> {
              tokenEndpoint
                  .accessTokenRequestConverter(accessTokenRequestConverter())
                  .errorResponseHandler(new WeCodingAuthenticationFailureEventHandler());
            })
        .clientAuthentication(
            clientAuthentication -> {
              clientAuthentication.errorResponseHandler(new WeCodingAuthenticationFailureEventHandler());
            })
        .authorizationService(authorizationService)
        .oidc(Customizer.withDefaults());

    DefaultSecurityFilterChain securityFilterChain =
        http.apply(new FormIdentityLoginConfigurer()).and().build();

    applyCustomOAuth2GrantAuthenticationProvider(http);

    return securityFilterChain;
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

  private AuthenticationConverter accessTokenRequestConverter() {
    return new DelegatingAuthenticationConverter(
        Arrays.asList(
            new OAuth2ResourceOwnerPasswordAuthenticationConverter(),
            new OAuth2RefreshTokenAuthenticationConverter(),
            new OAuth2ClientCredentialsAuthenticationConverter(),
            new OAuth2AuthorizationCodeAuthenticationConverter(),
            new OAuth2AuthorizationCodeRequestAuthenticationConverter(),
            new JwtClientAssertionAuthenticationConverter()));
  }

  @SuppressWarnings("unchecked")
  private void applyCustomOAuth2GrantAuthenticationProvider(HttpSecurity http) {
    AuthenticationManager authenticationManager = http.getSharedObject(AuthenticationManager.class);
    OAuth2AuthorizationService authorizationService =
        http.getSharedObject(OAuth2AuthorizationService.class);

    OAuth2TokenGenerator<? extends OAuth2Token> tokenGenerator =
        http.getSharedObject(OAuth2TokenGenerator.class);

    OAuth2ResourceOwnerPasswordAuthenticationProvider resourceOwnerPasswordAuthenticationProvider =
        new OAuth2ResourceOwnerPasswordAuthenticationProvider(
            authenticationManager, authorizationService, tokenGenerator);

    // 处理 UsernamePasswordAuthenticationToken
    http.authenticationProvider(weCodingDaoAuthenticationProvider);
    // 处理 OAuth2ResourceOwnerPasswordAuthenticationToken
    http.authenticationProvider(resourceOwnerPasswordAuthenticationProvider);
  }

  @Bean
  public AuthorizationServerSettings authorizationServerSettings() {
    return AuthorizationServerSettings.builder().build();
  }
}
