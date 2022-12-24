package top.wecoding.iam.server.security.configurers;

import java.net.URI;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import org.springframework.http.HttpStatus;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.config.annotation.web.configurers.ExceptionHandlingConfigurer;
import org.springframework.security.oauth2.server.authorization.config.annotation.web.configurers.OAuth2ClientAuthenticationConfigurer;
import org.springframework.security.oauth2.server.authorization.settings.AuthorizationServerSettings;
import org.springframework.security.web.authentication.HttpStatusEntryPoint;
import org.springframework.security.web.util.matcher.OrRequestMatcher;
import org.springframework.security.web.util.matcher.RequestMatcher;

/**
 * @author liuyuhui
 * @since 0.5
 */
public class WeCodingAuthorizationServerConfigurer
    extends AbstractHttpConfigurer<WeCodingAuthorizationServerConfigurer, HttpSecurity> {

  private final Map<Class<? extends AbstractLoginFilterConfigurer>, AbstractLoginFilterConfigurer>
      configurers = createConfigurers();
  private RequestMatcher endpointsMatcher;

  /**
   * Configures the WeCoding Password Login Endpoint.
   *
   * @param passwordLoginEndpointConfigurerCustomizer the {@link Customizer} providing access to the
   *     {@link Oauth2ResourceOwnerTokenEndpointFilterConfigurer}
   * @return the {@link WeCodingAuthorizationServerConfigurer} for further configuration
   */
  public WeCodingAuthorizationServerConfigurer passwordLoginEndpoint(
      Customizer<Oauth2ResourceOwnerTokenEndpointFilterConfigurer>
          passwordLoginEndpointConfigurerCustomizer) {
    passwordLoginEndpointConfigurerCustomizer.customize(
        getConfigurer(Oauth2ResourceOwnerTokenEndpointFilterConfigurer.class));
    return this;
  }

  public WeCodingAuthorizationServerConfigurer clientAuthentication(
      Customizer<OAuth2ClientAuthenticationConfigurer> clientAuthenticationCustomizer) {
    clientAuthenticationCustomizer.customize(
        getConfigurer(OAuth2ClientAuthenticationConfigurer.class));
    return this;
  }

  /**
   * Returns a {@link RequestMatcher} for the authorization server endpoints.
   *
   * @return a {@link RequestMatcher} for the authorization server endpoints
   */
  public RequestMatcher getEndpointsMatcher() {
    // Return a deferred RequestMatcher
    // since endpointsMatcher is constructed in init(HttpSecurity).
    return (request) -> this.endpointsMatcher.matches(request);
  }

  @Override
  public void init(HttpSecurity httpSecurity) {
    List<RequestMatcher> requestMatchers = new ArrayList<>();
    this.configurers
        .values()
        .forEach(
            configurer -> {
              configurer.init(httpSecurity);
              requestMatchers.add(configurer.getRequestMatcher());
            });
    this.endpointsMatcher = new OrRequestMatcher(requestMatchers);

    ExceptionHandlingConfigurer<HttpSecurity> exceptionHandling =
        httpSecurity.getConfigurer(ExceptionHandlingConfigurer.class);
    if (exceptionHandling != null) {
      exceptionHandling.defaultAuthenticationEntryPointFor(
          new HttpStatusEntryPoint(HttpStatus.UNAUTHORIZED),
          new OrRequestMatcher(
              getRequestMatcher(Oauth2ResourceOwnerTokenEndpointFilterConfigurer.class)));
    }
  }

  @Override
  public void configure(HttpSecurity httpSecurity) {
    this.configurers.values().forEach(configurer -> configurer.configure(httpSecurity));
  }

  private Map<Class<? extends AbstractLoginFilterConfigurer>, AbstractLoginFilterConfigurer>
      createConfigurers() {
    Map<Class<? extends AbstractLoginFilterConfigurer>, AbstractLoginFilterConfigurer> configurers =
        new LinkedHashMap<>();
    configurers.put(
        WeCodingClientAuthenticationConfigurer.class,
        new WeCodingClientAuthenticationConfigurer(this::postProcess));
    configurers.put(
        Oauth2ResourceOwnerTokenEndpointFilterConfigurer.class,
        new Oauth2ResourceOwnerTokenEndpointFilterConfigurer(this::postProcess));
    return configurers;
  }

  @SuppressWarnings("unchecked")
  private <T> T getConfigurer(Class<T> type) {
    return (T) this.configurers.get(type);
  }

  private <T extends AbstractLoginFilterConfigurer> void addConfigurer(
      Class<T> configurerType, T configurer) {
    this.configurers.put(configurerType, configurer);
  }

  private <T extends AbstractLoginFilterConfigurer> RequestMatcher getRequestMatcher(
      Class<T> configurerType) {
    T configurer = getConfigurer(configurerType);
    return configurer != null ? configurer.getRequestMatcher() : null;
  }

  private void validateAuthorizationServerSettings(
      AuthorizationServerSettings authorizationServerSettings) {
    if (authorizationServerSettings.getIssuer() != null) {
      URI issuerUri;
      try {
        issuerUri = new URI(authorizationServerSettings.getIssuer());
        issuerUri.toURL();
      } catch (Exception ex) {
        throw new IllegalArgumentException("issuer must be a valid URL", ex);
      }
      // rfc8414 https://datatracker.ietf.org/doc/html/rfc8414#section-2
      if (issuerUri.getQuery() != null || issuerUri.getFragment() != null) {
        throw new IllegalArgumentException("issuer cannot contain query or fragment component");
      }
    }
  }
}
