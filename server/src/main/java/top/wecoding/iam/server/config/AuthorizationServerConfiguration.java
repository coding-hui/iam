package top.wecoding.iam.server.config;

import lombok.RequiredArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configurers.oauth2.server.authorization.OAuth2AuthorizationServerConfigurer;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.security.oauth2.server.authorization.config.ProviderSettings;
import org.springframework.security.oauth2.server.authorization.token.DelegatingOAuth2TokenGenerator;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenGenerator;
import org.springframework.security.oauth2.server.authorization.web.authentication.*;
import org.springframework.security.web.DefaultSecurityFilterChain;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.AuthenticationConverter;
import org.springframework.security.web.util.matcher.RequestMatcher;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.iam.server.authentication.dao.IAMDaoAuthenticationProvider;
import top.wecoding.iam.server.authentication.handler.IAMAuthenticationFailureEventHandler;
import top.wecoding.iam.server.authentication.handler.IAMAuthenticationSuccessEventHandler;
import top.wecoding.iam.server.authentication.password.OAuth2ResourceOwnerPasswordAuthenticationConverter;
import top.wecoding.iam.server.authentication.password.OAuth2ResourceOwnerPasswordAuthenticationProvider;
import top.wecoding.iam.server.authentication.token.IAMOAuth2AccessTokenGenerator;
import top.wecoding.iam.server.authentication.token.IAMOAuth2RefreshTokenGenerator;
import top.wecoding.iam.server.authentication.token.IAMOAuth2TokenCustomizer;
import top.wecoding.iam.server.config.configurers.FormIdentityLoginConfigurer;

import java.util.Arrays;

import static top.wecoding.core.constant.SecurityConstants.CUSTOM_TOKEN_ENDPOINT_URI;
import static top.wecoding.core.constant.SecurityConstants.PROJECT_LICENSE;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Configuration
@RequiredArgsConstructor
public class AuthorizationServerConfiguration {

  private final OAuth2AuthorizationService authorizationService;

  @Bean
  @Order(Ordered.HIGHEST_PRECEDENCE)
  public SecurityFilterChain authorizationServerSecurityFilterChain(HttpSecurity http)
      throws Exception {
    OAuth2AuthorizationServerConfigurer<HttpSecurity> authorizationServerConfigurer =
        new OAuth2AuthorizationServerConfigurer<>();

    RequestMatcher endpointsMatcher = authorizationServerConfigurer.getEndpointsMatcher();

    http.apply(
        authorizationServerConfigurer
            .tokenEndpoint(
                (tokenEndpoint) ->
                    tokenEndpoint
                        .accessTokenRequestConverter(accessTokenRequestConverter())
                        .accessTokenResponseHandler(new IAMAuthenticationSuccessEventHandler())
                        .errorResponseHandler(new IAMAuthenticationFailureEventHandler()))
            .clientAuthentication(
                oAuth2ClientAuthenticationConfigurer ->
                    oAuth2ClientAuthenticationConfigurer.errorResponseHandler(
                        new IAMAuthenticationFailureEventHandler()))
            .authorizationEndpoint(
                authorizationEndpoint ->
                    authorizationEndpoint.consentPage(SecurityConstants.CUSTOM_CONSENT_PAGE_URI)));

    DefaultSecurityFilterChain securityFilterChain =
        http.requestMatcher(endpointsMatcher)
            .authorizeRequests(authorizeRequests -> authorizeRequests.anyRequest().authenticated())
            .apply(
                authorizationServerConfigurer
                    .authorizationService(authorizationService)
                    .providerSettings(
                        ProviderSettings.builder()
                            .tokenEndpoint(CUSTOM_TOKEN_ENDPOINT_URI)
                            .issuer(PROJECT_LICENSE)
                            .build()))
            .and()
            .apply(new FormIdentityLoginConfigurer())
            .and()
            .build();

    addCustomOAuth2GrantAuthenticationProvider(http);
    return securityFilterChain;
  }

  @Bean
  @SuppressWarnings("rawtypes")
  public OAuth2TokenGenerator oAuth2TokenGenerator() {
    IAMOAuth2AccessTokenGenerator accessTokenGenerator = new IAMOAuth2AccessTokenGenerator();
    accessTokenGenerator.setAccessTokenCustomizer(new IAMOAuth2TokenCustomizer());
    return new DelegatingOAuth2TokenGenerator(
        accessTokenGenerator, new IAMOAuth2RefreshTokenGenerator());
  }

  /**
   * request -> xToken 注入请求转换器
   *
   * @return DelegatingAuthenticationConverter
   */
  private AuthenticationConverter accessTokenRequestConverter() {
    return new DelegatingAuthenticationConverter(
        Arrays.asList(
            new OAuth2ResourceOwnerPasswordAuthenticationConverter(),
            // new OAuth2ResourceOwnerSmsAuthenticationConverter(),
            new OAuth2RefreshTokenAuthenticationConverter(),
            new OAuth2ClientCredentialsAuthenticationConverter(),
            new OAuth2AuthorizationCodeAuthenticationConverter(),
            new OAuth2AuthorizationCodeRequestAuthenticationConverter()));
  }

  @SuppressWarnings("unchecked")
  private void addCustomOAuth2GrantAuthenticationProvider(HttpSecurity http) {
    AuthenticationManager authenticationManager = http.getSharedObject(AuthenticationManager.class);
    OAuth2AuthorizationService authorizationService =
        http.getSharedObject(OAuth2AuthorizationService.class);

    OAuth2ResourceOwnerPasswordAuthenticationProvider resourceOwnerPasswordAuthenticationProvider =
        new OAuth2ResourceOwnerPasswordAuthenticationProvider(
            authenticationManager, authorizationService, oAuth2TokenGenerator());

    // 处理 UsernamePasswordAuthenticationToken
    http.authenticationProvider(new IAMDaoAuthenticationProvider());
    // 处理 OAuth2ResourceOwnerPasswordAuthenticationToken
    http.authenticationProvider(resourceOwnerPasswordAuthenticationProvider);
  }
}
