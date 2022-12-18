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
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configurers.oauth2.server.resource.OAuth2ResourceServerConfigurer;
import org.springframework.security.oauth2.jwt.JwtDecoder;
import org.springframework.security.oauth2.jwt.NimbusJwtEncoder;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.security.oauth2.server.authorization.config.annotation.web.configuration.OAuth2AuthorizationServerConfiguration;
import org.springframework.security.oauth2.server.authorization.config.annotation.web.configurers.OAuth2AuthorizationServerConfigurer;
import org.springframework.security.oauth2.server.authorization.config.annotation.web.configurers.OAuth2ClientAuthenticationConfigurer;
import org.springframework.security.oauth2.server.authorization.config.annotation.web.configurers.OAuth2TokenEndpointConfigurer;
import org.springframework.security.oauth2.server.authorization.settings.AuthorizationServerSettings;
import org.springframework.security.oauth2.server.authorization.token.*;
import org.springframework.security.oauth2.server.authorization.web.authentication.*;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.AuthenticationConverter;
import org.springframework.security.web.util.matcher.OrRequestMatcher;
import org.springframework.security.web.util.matcher.RequestMatcher;
import top.wecoding.iam.framework.security.handler.WeCodingAuthenticationFailureEventHandler;
import top.wecoding.iam.framework.security.jose.Jwks;
import top.wecoding.iam.framework.security.web.ResourceAuthExceptionEntryPoint;
import top.wecoding.iam.server.security.authorization.authentication.password.OAuth2ResourceOwnerPasswordAuthenticationConverter;
import top.wecoding.iam.server.security.authorization.token.IAMOAuth2TokenCustomizer;
import top.wecoding.iam.server.security.configurers.FormIdentityLoginConfigurer;
import top.wecoding.iam.server.security.configurers.WeCodingAuthorizationServerConfigurer;
import top.wecoding.iam.server.security.handler.SsoAuthenticationSuccessHandler;

import java.util.Arrays;

import static top.wecoding.iam.common.constant.WeCodingSettingNames.AuthorizationServer.RESOURCE_OWNER_TOKEN_ENDPOINT;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
@RequiredArgsConstructor
@Configuration(proxyBeanMethods = false)
public class AuthorizationServerConfig {

  protected final ResourceAuthExceptionEntryPoint resourceAuthExceptionEntryPoint;

  private final OAuth2AuthorizationService authorizationService;

  @Bean
  @Order(Ordered.HIGHEST_PRECEDENCE)
  public SecurityFilterChain authorizationServerSecurityFilterChain(
      HttpSecurity http,
      Customizer<OAuth2TokenEndpointConfigurer> tokenEndpointCustomizer,
      Customizer<OAuth2ClientAuthenticationConfigurer> clientAuthenticationCustomizer,
      Customizer<OAuth2ResourceServerConfigurer<HttpSecurity>> oauth2ResourceServerCustomizer)
      throws Exception {
    OAuth2AuthorizationServerConfigurer authorizationServerConfigurer =
        new OAuth2AuthorizationServerConfigurer();
    WeCodingAuthorizationServerConfigurer weCodingAuthorizationServerConfigurer =
        new WeCodingAuthorizationServerConfigurer();

    RequestMatcher requestMatcher =
        new OrRequestMatcher(
            weCodingAuthorizationServerConfigurer.getEndpointsMatcher(),
            authorizationServerConfigurer.getEndpointsMatcher());

    http.securityMatchers(matchers -> matchers.requestMatchers(requestMatcher))
        .authorizeHttpRequests(authorize -> authorize.anyRequest().authenticated())
        .csrf(csrf -> csrf.ignoringRequestMatchers(requestMatcher));

    http.apply(authorizationServerConfigurer)
        .tokenEndpoint(tokenEndpointCustomizer)
        .clientAuthentication(clientAuthenticationCustomizer)
        .authorizationService(authorizationService)
        .oidc(Customizer.withDefaults());

    http.oauth2ResourceServer(oauth2ResourceServerCustomizer);

    http.apply(weCodingAuthorizationServerConfigurer)
        .passwordLoginEndpoint(
            passwordLogin ->
                passwordLogin
                    .accessTokenRequestConverter(accessTokenRequestConverter())
                    .errorResponseHandler(new WeCodingAuthenticationFailureEventHandler()));

    http.apply(new FormIdentityLoginConfigurer());

    return http.build();
  }

  @Bean
  public Customizer<OAuth2TokenEndpointConfigurer> tokenEndpointCustomizer() {
    return customizer ->
        customizer
            .accessTokenRequestConverter(accessTokenRequestConverter())
            .accessTokenResponseHandler(new SsoAuthenticationSuccessHandler())
            .errorResponseHandler(new WeCodingAuthenticationFailureEventHandler());
  }

  @Bean
  public Customizer<OAuth2ClientAuthenticationConfigurer> clientAuthenticationCustomizer() {
    return customizer ->
        customizer.errorResponseHandler(new WeCodingAuthenticationFailureEventHandler());
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

  @Bean
  @SuppressWarnings("rawtypes")
  public OAuth2TokenGenerator oAuth2TokenGenerator(JWKSource<SecurityContext> jwkSource) {
    OAuth2AccessTokenGenerator accessTokenGenerator = new OAuth2AccessTokenGenerator();
    accessTokenGenerator.setAccessTokenCustomizer(new IAMOAuth2TokenCustomizer());
    JwtGenerator jwtGenerator = new JwtGenerator(new NimbusJwtEncoder(jwkSource));
    return new DelegatingOAuth2TokenGenerator(
        accessTokenGenerator, new OAuth2RefreshTokenGenerator(), jwtGenerator);
  }

  @Bean
  public AuthorizationServerSettings authorizationServerSettings() {
    return AuthorizationServerSettings.builder()
        .authorizationEndpoint("/oauth2/authorize")
        .tokenEndpoint("/oauth2/token")
        .jwkSetEndpoint("/oauth2/jwks")
        .tokenRevocationEndpoint("/oauth2/revoke")
        .tokenIntrospectionEndpoint("/oauth2/introspect")
        .oidcClientRegistrationEndpoint("/connect/register")
        .oidcUserInfoEndpoint("/userinfo")
        .setting(RESOURCE_OWNER_TOKEN_ENDPOINT, "/api/v1/signin")
        .build();
  }

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
}
