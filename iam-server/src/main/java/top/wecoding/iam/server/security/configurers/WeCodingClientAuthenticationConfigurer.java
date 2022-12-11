package top.wecoding.iam.server.security.configurers;

import jakarta.servlet.http.HttpServletRequest;
import org.springframework.http.HttpMethod;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.AuthenticationProvider;
import org.springframework.security.config.annotation.ObjectPostProcessor;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.core.context.SecurityContext;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.oauth2.core.OAuth2Error;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.security.oauth2.server.authorization.authentication.ClientSecretAuthenticationProvider;
import org.springframework.security.oauth2.server.authorization.authentication.JwtClientAssertionAuthenticationProvider;
import org.springframework.security.oauth2.server.authorization.authentication.OAuth2ClientAuthenticationToken;
import org.springframework.security.oauth2.server.authorization.authentication.PublicClientAuthenticationProvider;
import org.springframework.security.oauth2.server.authorization.client.RegisteredClientRepository;
import org.springframework.security.oauth2.server.authorization.settings.AuthorizationServerSettings;
import org.springframework.security.oauth2.server.authorization.web.authentication.*;
import org.springframework.security.web.authentication.AuthenticationConverter;
import org.springframework.security.web.authentication.AuthenticationFailureHandler;
import org.springframework.security.web.authentication.AuthenticationSuccessHandler;
import org.springframework.security.web.authentication.preauth.AbstractPreAuthenticatedProcessingFilter;
import org.springframework.security.web.util.matcher.AntPathRequestMatcher;
import org.springframework.security.web.util.matcher.OrRequestMatcher;
import org.springframework.security.web.util.matcher.RequestMatcher;
import org.springframework.util.Assert;
import top.wecoding.iam.common.constant.WeCodingSettingNames;
import top.wecoding.iam.server.security.web.WeCodingClientAuthenticationFilter;

import java.util.ArrayList;
import java.util.List;
import java.util.function.Consumer;

/**
 * @author liuyuhui
 * @since 0.5
 */
public class WeCodingClientAuthenticationConfigurer extends AbstractLoginFilterConfigurer {

  private final List<AuthenticationConverter> authenticationConverters = new ArrayList<>();
  private final List<AuthenticationProvider> authenticationProviders = new ArrayList<>();
  private RequestMatcher requestMatcher;
  private Consumer<List<AuthenticationConverter>> authenticationConvertersConsumer =
      (authenticationConverters) -> {};
  private Consumer<List<AuthenticationProvider>> authenticationProvidersConsumer =
      (authenticationProviders) -> {};
  private AuthenticationSuccessHandler authenticationSuccessHandler;
  private AuthenticationFailureHandler errorResponseHandler;

  /** Restrict for internal use only. */
  WeCodingClientAuthenticationConfigurer(ObjectPostProcessor<Object> objectPostProcessor) {
    super(objectPostProcessor);
  }

  private static List<AuthenticationConverter> createDefaultAuthenticationConverters() {
    List<AuthenticationConverter> authenticationConverters = new ArrayList<>();

    authenticationConverters.add(new JwtClientAssertionAuthenticationConverter());
    authenticationConverters.add(new ClientSecretBasicAuthenticationConverter());
    authenticationConverters.add(new ClientSecretPostAuthenticationConverter());
    authenticationConverters.add(new PublicClientAuthenticationConverter());

    return authenticationConverters;
  }

  private static List<AuthenticationProvider> createDefaultAuthenticationProviders(
      HttpSecurity httpSecurity) {
    List<AuthenticationProvider> authenticationProviders = new ArrayList<>();

    RegisteredClientRepository registeredClientRepository =
        OAuth2ConfigurerUtils.getRegisteredClientRepository(httpSecurity);
    OAuth2AuthorizationService authorizationService =
        OAuth2ConfigurerUtils.getAuthorizationService(httpSecurity);

    JwtClientAssertionAuthenticationProvider jwtClientAssertionAuthenticationProvider =
        new JwtClientAssertionAuthenticationProvider(
            registeredClientRepository, authorizationService);
    authenticationProviders.add(jwtClientAssertionAuthenticationProvider);

    ClientSecretAuthenticationProvider clientSecretAuthenticationProvider =
        new ClientSecretAuthenticationProvider(registeredClientRepository, authorizationService);
    PasswordEncoder passwordEncoder =
        OAuth2ConfigurerUtils.getOptionalBean(httpSecurity, PasswordEncoder.class);
    if (passwordEncoder != null) {
      clientSecretAuthenticationProvider.setPasswordEncoder(passwordEncoder);
    }
    authenticationProviders.add(clientSecretAuthenticationProvider);

    PublicClientAuthenticationProvider publicClientAuthenticationProvider =
        new PublicClientAuthenticationProvider(registeredClientRepository, authorizationService);
    authenticationProviders.add(publicClientAuthenticationProvider);

    return authenticationProviders;
  }

  /**
   * Adds an {@link AuthenticationConverter} used when attempting to extract client credentials from
   * {@link HttpServletRequest} to an instance of {@link OAuth2ClientAuthenticationToken} used for
   * authenticating the client.
   *
   * @param authenticationConverter an {@link AuthenticationConverter} used when attempting to
   *     extract client credentials from {@link HttpServletRequest}
   * @return the {@link WeCodingClientAuthenticationConfigurer} for further configuration
   */
  public WeCodingClientAuthenticationConfigurer authenticationConverter(
      AuthenticationConverter authenticationConverter) {
    Assert.notNull(authenticationConverter, "authenticationConverter cannot be null");
    this.authenticationConverters.add(authenticationConverter);
    return this;
  }

  /**
   * Sets the {@code Consumer} providing access to the {@code List} of default and (optionally)
   * added {@link #authenticationConverter(AuthenticationConverter) AuthenticationConverter}'s
   * allowing the ability to add, remove, or customize a specific {@link AuthenticationConverter}.
   *
   * @param authenticationConvertersConsumer the {@code Consumer} providing access to the {@code
   *     List} of default and (optionally) added {@link AuthenticationConverter}'s
   * @return the {@link WeCodingClientAuthenticationConfigurer} for further configuration
   */
  public WeCodingClientAuthenticationConfigurer authenticationConverters(
      Consumer<List<AuthenticationConverter>> authenticationConvertersConsumer) {
    Assert.notNull(
        authenticationConvertersConsumer, "authenticationConvertersConsumer cannot be null");
    this.authenticationConvertersConsumer = authenticationConvertersConsumer;
    return this;
  }

  /**
   * Adds an {@link AuthenticationProvider} used for authenticating an {@link
   * OAuth2ClientAuthenticationToken}.
   *
   * @param authenticationProvider an {@link AuthenticationProvider} used for authenticating an
   *     {@link OAuth2ClientAuthenticationToken}
   * @return the {@link WeCodingClientAuthenticationConfigurer} for further configuration
   */
  public WeCodingClientAuthenticationConfigurer authenticationProvider(
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
   * @return the {@link WeCodingClientAuthenticationConfigurer} for further configuration
   */
  public WeCodingClientAuthenticationConfigurer authenticationProviders(
      Consumer<List<AuthenticationProvider>> authenticationProvidersConsumer) {
    Assert.notNull(
        authenticationProvidersConsumer, "authenticationProvidersConsumer cannot be null");
    this.authenticationProvidersConsumer = authenticationProvidersConsumer;
    return this;
  }

  /**
   * Sets the {@link AuthenticationSuccessHandler} used for handling a successful client
   * authentication and associating the {@link OAuth2ClientAuthenticationToken} to the {@link
   * SecurityContext}.
   *
   * @param authenticationSuccessHandler the {@link AuthenticationSuccessHandler} used for handling
   *     a successful client authentication
   * @return the {@link WeCodingClientAuthenticationConfigurer} for further configuration
   */
  public WeCodingClientAuthenticationConfigurer authenticationSuccessHandler(
      AuthenticationSuccessHandler authenticationSuccessHandler) {
    this.authenticationSuccessHandler = authenticationSuccessHandler;
    return this;
  }

  /**
   * Sets the {@link AuthenticationFailureHandler} used for handling a failed client authentication
   * and returning the {@link OAuth2Error Error Response}.
   *
   * @param errorResponseHandler the {@link AuthenticationFailureHandler} used for handling a failed
   *     client authentication
   * @return the {@link WeCodingClientAuthenticationConfigurer} for further configuration
   */
  public WeCodingClientAuthenticationConfigurer errorResponseHandler(
      AuthenticationFailureHandler errorResponseHandler) {
    this.errorResponseHandler = errorResponseHandler;
    return this;
  }

  @Override
  void init(HttpSecurity httpSecurity) {
    AuthorizationServerSettings authorizationServerSettings =
        OAuth2ConfigurerUtils.getAuthorizationServerSettings(httpSecurity);
    this.requestMatcher =
        new OrRequestMatcher(
            new AntPathRequestMatcher(
                authorizationServerSettings.getSetting(
                    WeCodingSettingNames.AuthorizationServer.RESOURCE_OWNER_TOKEN_ENDPOINT),
                HttpMethod.POST.name()));

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
    WeCodingClientAuthenticationFilter clientAuthenticationFilter =
        new WeCodingClientAuthenticationFilter(authenticationManager, this.requestMatcher);
    List<AuthenticationConverter> authenticationConverters =
        createDefaultAuthenticationConverters();
    if (!this.authenticationConverters.isEmpty()) {
      authenticationConverters.addAll(0, this.authenticationConverters);
    }
    this.authenticationConvertersConsumer.accept(authenticationConverters);
    clientAuthenticationFilter.setAuthenticationConverter(
        new DelegatingAuthenticationConverter(authenticationConverters));
    if (this.authenticationSuccessHandler != null) {
      clientAuthenticationFilter.setAuthenticationSuccessHandler(this.authenticationSuccessHandler);
    }
    if (this.errorResponseHandler != null) {
      clientAuthenticationFilter.setAuthenticationFailureHandler(this.errorResponseHandler);
    }
    httpSecurity.addFilterAfter(
        postProcess(clientAuthenticationFilter), AbstractPreAuthenticatedProcessingFilter.class);
  }

  @Override
  RequestMatcher getRequestMatcher() {
    return this.requestMatcher;
  }
}
