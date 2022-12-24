package top.wecoding.iam.server.security.configurers;

import jakarta.servlet.http.HttpServletRequest;
import java.util.ArrayList;
import java.util.List;
import java.util.function.Consumer;
import org.springframework.http.HttpMethod;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.AuthenticationProvider;
import org.springframework.security.config.annotation.ObjectPostProcessor;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.OAuth2Error;
import org.springframework.security.oauth2.core.OAuth2Token;
import org.springframework.security.oauth2.core.endpoint.OAuth2AccessTokenResponse;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.security.oauth2.server.authorization.authentication.*;
import org.springframework.security.oauth2.server.authorization.settings.AuthorizationServerSettings;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenGenerator;
import org.springframework.security.oauth2.server.authorization.web.authentication.DelegatingAuthenticationConverter;
import org.springframework.security.oauth2.server.authorization.web.authentication.OAuth2AuthorizationCodeAuthenticationConverter;
import org.springframework.security.oauth2.server.authorization.web.authentication.OAuth2ClientCredentialsAuthenticationConverter;
import org.springframework.security.oauth2.server.authorization.web.authentication.OAuth2RefreshTokenAuthenticationConverter;
import org.springframework.security.web.access.intercept.AuthorizationFilter;
import org.springframework.security.web.authentication.AuthenticationConverter;
import org.springframework.security.web.authentication.AuthenticationFailureHandler;
import org.springframework.security.web.authentication.AuthenticationSuccessHandler;
import org.springframework.security.web.util.matcher.AntPathRequestMatcher;
import org.springframework.security.web.util.matcher.RequestMatcher;
import org.springframework.util.Assert;
import top.wecoding.iam.common.constant.WeCodingSettingNames;
import top.wecoding.iam.server.security.authorization.authentication.OAuth2ResourceOwnerBaseAuthenticationToken;
import top.wecoding.iam.server.security.authorization.authentication.WeCodingDaoAuthenticationProvider;
import top.wecoding.iam.server.security.authorization.authentication.password.OAuth2ResourceOwnerPasswordAuthenticationProvider;
import top.wecoding.iam.server.security.web.Oauth2ResourceOwnerTokenEndpointFilter;

/**
 * Configurer for the WeCoding Password Login Endpoint.
 *
 * @author liuyuhui
 * @since 0.5
 * @see WeCodingAuthorizationServerConfigurer#passwordLoginEndpoint
 * @see Oauth2ResourceOwnerTokenEndpointFilter
 */
public final class Oauth2ResourceOwnerTokenEndpointFilterConfigurer
    extends AbstractLoginFilterConfigurer {

  private final List<AuthenticationConverter> accessTokenRequestConverters = new ArrayList<>();
  private final List<AuthenticationProvider> authenticationProviders = new ArrayList<>();
  private RequestMatcher requestMatcher;
  private Consumer<List<AuthenticationConverter>> accessTokenRequestConvertersConsumer =
      (accessTokenRequestConverters) -> {};
  private Consumer<List<AuthenticationProvider>> authenticationProvidersConsumer =
      (authenticationProviders) -> {};
  private AuthenticationSuccessHandler accessTokenResponseHandler;
  private AuthenticationFailureHandler errorResponseHandler;

  /** Restrict for internal use only. */
  Oauth2ResourceOwnerTokenEndpointFilterConfigurer(
      ObjectPostProcessor<Object> objectPostProcessor) {
    super(objectPostProcessor);
  }

  /**
   * Adds an {@link AuthenticationConverter} used when attempting to extract an Access Token Request
   * from {@link HttpServletRequest} to an instance of {@link
   * OAuth2ResourceOwnerBaseAuthenticationToken} used for authenticating the authorization grant.
   *
   * @param accessTokenRequestConverter an {@link AuthenticationConverter} used when attempting to
   *     extract an Access Token Request from {@link HttpServletRequest}
   * @return the {@link Oauth2ResourceOwnerTokenEndpointFilterConfigurer} for further configuration
   */
  public Oauth2ResourceOwnerTokenEndpointFilterConfigurer accessTokenRequestConverter(
      AuthenticationConverter accessTokenRequestConverter) {
    Assert.notNull(accessTokenRequestConverter, "accessTokenRequestConverter cannot be null");
    this.accessTokenRequestConverters.add(accessTokenRequestConverter);
    return this;
  }

  /**
   * Sets the {@code Consumer} providing access to the {@code List} of default and (optionally)
   * added {@link #accessTokenRequestConverter(AuthenticationConverter) AuthenticationConverter}'s
   * allowing the ability to add, remove, or customize a specific {@link AuthenticationConverter}.
   *
   * @param accessTokenRequestConvertersConsumer the {@code Consumer} providing access to the {@code
   *     List} of default and (optionally) added {@link AuthenticationConverter}'s
   * @return the {@link Oauth2ResourceOwnerTokenEndpointFilterConfigurer} for further configuration
   * @since 0.4.0
   */
  public Oauth2ResourceOwnerTokenEndpointFilterConfigurer accessTokenRequestConverters(
      Consumer<List<AuthenticationConverter>> accessTokenRequestConvertersConsumer) {
    Assert.notNull(
        accessTokenRequestConvertersConsumer,
        "accessTokenRequestConvertersConsumer cannot be null");
    this.accessTokenRequestConvertersConsumer = accessTokenRequestConvertersConsumer;
    return this;
  }

  /**
   * Adds an {@link AuthenticationProvider} used for authenticating a type of {@link
   * OAuth2AuthorizationGrantAuthenticationToken}.
   *
   * @param authenticationProvider an {@link AuthenticationProvider} used for authenticating a type
   *     of {@link OAuth2AuthorizationGrantAuthenticationToken}
   * @return the {@link Oauth2ResourceOwnerTokenEndpointFilterConfigurer} for further configuration
   */
  public Oauth2ResourceOwnerTokenEndpointFilterConfigurer authenticationProvider(
      AuthenticationProvider authenticationProvider) {
    Assert.notNull(authenticationProvider, "authenticationProvider cannot be null");
    this.authenticationProviders.add(authenticationProvider);
    return this;
  }

  /**
   * Sets the {@code Consumer} providing access to the {@code List} of default and (optionally)
   * added {@link #authenticationProvider(AuthenticationProvider) AuthenticationProvider}'s allowing
   * the ability to add, remove, or customize a specific {@link AuthenticationProvider}.
   *
   * @param authenticationProvidersConsumer the {@code Consumer} providing access to the {@code
   *     List} of default and (optionally) added {@link AuthenticationProvider}'s
   * @return the {@link Oauth2ResourceOwnerTokenEndpointFilterConfigurer} for further configuration
   * @since 0.4.0
   */
  public Oauth2ResourceOwnerTokenEndpointFilterConfigurer authenticationProviders(
      Consumer<List<AuthenticationProvider>> authenticationProvidersConsumer) {
    Assert.notNull(
        authenticationProvidersConsumer, "authenticationProvidersConsumer cannot be null");
    this.authenticationProvidersConsumer = authenticationProvidersConsumer;
    return this;
  }

  /**
   * Sets the {@link AuthenticationSuccessHandler} used for handling an {@link
   * OAuth2AccessTokenAuthenticationToken} and returning the {@link OAuth2AccessTokenResponse Access
   * Token Response}.
   *
   * @param accessTokenResponseHandler the {@link AuthenticationSuccessHandler} used for handling an
   *     {@link OAuth2AccessTokenAuthenticationToken}
   * @return the {@link Oauth2ResourceOwnerTokenEndpointFilterConfigurer} for further configuration
   */
  public Oauth2ResourceOwnerTokenEndpointFilterConfigurer accessTokenResponseHandler(
      AuthenticationSuccessHandler accessTokenResponseHandler) {
    this.accessTokenResponseHandler = accessTokenResponseHandler;
    return this;
  }

  /**
   * Sets the {@link AuthenticationFailureHandler} used for handling an {@link
   * OAuth2AuthenticationException} and returning the {@link OAuth2Error Error Response}.
   *
   * @param errorResponseHandler the {@link AuthenticationFailureHandler} used for handling an
   *     {@link OAuth2AuthenticationException}
   * @return the {@link Oauth2ResourceOwnerTokenEndpointFilterConfigurer} for further configuration
   */
  public Oauth2ResourceOwnerTokenEndpointFilterConfigurer errorResponseHandler(
      AuthenticationFailureHandler errorResponseHandler) {
    this.errorResponseHandler = errorResponseHandler;
    return this;
  }

  @Override
  void init(HttpSecurity httpSecurity) {
    AuthorizationServerSettings authorizationServerSettings =
        OAuth2ConfigurerUtils.getAuthorizationServerSettings(httpSecurity);
    this.requestMatcher =
        new AntPathRequestMatcher(
            authorizationServerSettings.getSetting(
                WeCodingSettingNames.AuthorizationServer.RESOURCE_OWNER_TOKEN_ENDPOINT),
            HttpMethod.POST.name());

    List<AuthenticationProvider> authenticationProviders =
        createDefaultAuthenticationProviders(httpSecurity);
    if (!this.authenticationProviders.isEmpty()) {
      authenticationProviders.addAll(0, this.authenticationProviders);
    }
    this.authenticationProvidersConsumer.accept(authenticationProviders);
    authenticationProviders.forEach(
        authenticationProvider ->
            httpSecurity.authenticationProvider(postProcess(authenticationProvider)));
  }

  @Override
  void configure(HttpSecurity httpSecurity) {
    AuthenticationManager authenticationManager =
        httpSecurity.getSharedObject(AuthenticationManager.class);
    OAuth2AuthorizationService authorizationService =
        OAuth2ConfigurerUtils.getAuthorizationService(httpSecurity);
    OAuth2TokenGenerator<? extends OAuth2Token> tokenGenerator =
        OAuth2ConfigurerUtils.getTokenGenerator(httpSecurity);
    AuthorizationServerSettings authorizationServerSettings =
        OAuth2ConfigurerUtils.getAuthorizationServerSettings(httpSecurity);

    Oauth2ResourceOwnerTokenEndpointFilter oauth2ResourceOwnerTokenEndpointFilter =
        new Oauth2ResourceOwnerTokenEndpointFilter(
            authenticationManager,
            authorizationServerSettings.getSetting(
                WeCodingSettingNames.AuthorizationServer.RESOURCE_OWNER_TOKEN_ENDPOINT));
    List<AuthenticationConverter> authenticationConverters =
        createDefaultAuthenticationConverters();
    if (!this.accessTokenRequestConverters.isEmpty()) {
      authenticationConverters.addAll(0, this.accessTokenRequestConverters);
    }
    this.accessTokenRequestConvertersConsumer.accept(authenticationConverters);
    oauth2ResourceOwnerTokenEndpointFilter.setAuthenticationConverter(
        new DelegatingAuthenticationConverter(authenticationConverters));
    if (this.accessTokenResponseHandler != null) {
      oauth2ResourceOwnerTokenEndpointFilter.setAuthenticationSuccessHandler(
          this.accessTokenResponseHandler);
    }
    if (this.errorResponseHandler != null) {
      oauth2ResourceOwnerTokenEndpointFilter.setAuthenticationFailureHandler(
          this.errorResponseHandler);
    }

    OAuth2ResourceOwnerPasswordAuthenticationProvider resourceOwnerPasswordAuthenticationProvider =
        new OAuth2ResourceOwnerPasswordAuthenticationProvider(
            authenticationManager, authorizationService, tokenGenerator);
    authenticationProviders.add(resourceOwnerPasswordAuthenticationProvider);
    httpSecurity.authenticationProvider(
        getObjectPostProcessor().postProcess(resourceOwnerPasswordAuthenticationProvider));

    httpSecurity.addFilterAfter(
        postProcess(oauth2ResourceOwnerTokenEndpointFilter), AuthorizationFilter.class);
  }

  @Override
  RequestMatcher getRequestMatcher() {
    return this.requestMatcher;
  }

  private List<AuthenticationConverter> createDefaultAuthenticationConverters() {
    List<AuthenticationConverter> authenticationConverters = new ArrayList<>();

    authenticationConverters.add(new OAuth2AuthorizationCodeAuthenticationConverter());
    authenticationConverters.add(new OAuth2RefreshTokenAuthenticationConverter());
    authenticationConverters.add(new OAuth2ClientCredentialsAuthenticationConverter());

    return authenticationConverters;
  }

  private List<AuthenticationProvider> createDefaultAuthenticationProviders(
      HttpSecurity httpSecurity) {
    List<AuthenticationProvider> authenticationProviders = new ArrayList<>();

    OAuth2AuthorizationService authorizationService =
        OAuth2ConfigurerUtils.getAuthorizationService(httpSecurity);
    OAuth2TokenGenerator<? extends OAuth2Token> tokenGenerator =
        OAuth2ConfigurerUtils.getTokenGenerator(httpSecurity);
    UserDetailsService userDetailsService =
        OAuth2ConfigurerUtils.getBean(httpSecurity, UserDetailsService.class);

    OAuth2AuthorizationCodeAuthenticationProvider authorizationCodeAuthenticationProvider =
        new OAuth2AuthorizationCodeAuthenticationProvider(authorizationService, tokenGenerator);
    authenticationProviders.add(authorizationCodeAuthenticationProvider);

    OAuth2RefreshTokenAuthenticationProvider refreshTokenAuthenticationProvider =
        new OAuth2RefreshTokenAuthenticationProvider(authorizationService, tokenGenerator);
    authenticationProviders.add(refreshTokenAuthenticationProvider);

    OAuth2ClientCredentialsAuthenticationProvider clientCredentialsAuthenticationProvider =
        new OAuth2ClientCredentialsAuthenticationProvider(authorizationService, tokenGenerator);
    authenticationProviders.add(clientCredentialsAuthenticationProvider);

    WeCodingDaoAuthenticationProvider weCodingDaoAuthenticationProvider =
        new WeCodingDaoAuthenticationProvider(userDetailsService);
    authenticationProviders.add(weCodingDaoAuthenticationProvider);

    return authenticationProviders;
  }
}
