package top.wecoding.iam.server.security.web;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.time.temporal.ChronoUnit;
import java.util.Arrays;
import java.util.Map;
import lombok.extern.slf4j.Slf4j;
import org.springframework.core.log.LogMessage;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.converter.HttpMessageConverter;
import org.springframework.http.server.ServletServerHttpResponse;
import org.springframework.security.authentication.AbstractAuthenticationToken;
import org.springframework.security.authentication.AuthenticationDetailsSource;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.core.*;
import org.springframework.security.oauth2.core.endpoint.OAuth2AccessTokenResponse;
import org.springframework.security.oauth2.core.http.converter.OAuth2AccessTokenResponseHttpMessageConverter;
import org.springframework.security.oauth2.core.http.converter.OAuth2ErrorHttpMessageConverter;
import org.springframework.security.oauth2.server.authorization.authentication.OAuth2AccessTokenAuthenticationToken;
import org.springframework.security.oauth2.server.authorization.web.authentication.DelegatingAuthenticationConverter;
import org.springframework.security.oauth2.server.authorization.web.authentication.OAuth2AuthorizationCodeAuthenticationConverter;
import org.springframework.security.oauth2.server.authorization.web.authentication.OAuth2ClientCredentialsAuthenticationConverter;
import org.springframework.security.oauth2.server.authorization.web.authentication.OAuth2RefreshTokenAuthenticationConverter;
import org.springframework.security.web.authentication.AuthenticationConverter;
import org.springframework.security.web.authentication.AuthenticationFailureHandler;
import org.springframework.security.web.authentication.AuthenticationSuccessHandler;
import org.springframework.security.web.authentication.WebAuthenticationDetailsSource;
import org.springframework.security.web.util.matcher.AntPathRequestMatcher;
import org.springframework.security.web.util.matcher.RequestMatcher;
import org.springframework.util.Assert;
import org.springframework.util.CollectionUtils;
import org.springframework.web.filter.OncePerRequestFilter;
import top.wecoding.iam.common.constant.AuthParameterNames;
import top.wecoding.iam.common.constant.SecurityConstants;
import top.wecoding.iam.common.convert.RestAccessTokenResponseHttpMessageConverter;
import top.wecoding.iam.common.convert.RestOAuth2ErrorParametersConverter;
import top.wecoding.iam.server.security.authorization.authentication.OAuth2ResourceOwnerBaseAuthenticationToken;
import top.wecoding.iam.server.security.configurers.Oauth2ResourceOwnerTokenEndpointFilterConfigurer;
import top.wecoding.iam.server.util.LogUtil;

/**
 * @author liuyuhui
 * @since 0.5
 * @see Oauth2ResourceOwnerTokenEndpointFilterConfigurer
 */
@Slf4j
public class Oauth2ResourceOwnerTokenEndpointFilter extends OncePerRequestFilter {

  private static final String DEFAULT_TOKEN_ENDPOINT_URI = "/oauth2/token";

  private static final String DEFAULT_ERROR_URI =
      "https://datatracker.ietf.org/doc/html/rfc6749#section-5.2";
  private final AuthenticationManager authenticationManager;
  private final RequestMatcher tokenEndpointMatcher;
  private final HttpMessageConverter<OAuth2AccessTokenResponse> accessTokenHttpResponseConverter;
  private final HttpMessageConverter<OAuth2Error> errorHttpResponseConverter;
  private AuthenticationDetailsSource<HttpServletRequest, ?> authenticationDetailsSource =
      new WebAuthenticationDetailsSource();
  private AuthenticationConverter authenticationConverter;
  private AuthenticationSuccessHandler authenticationSuccessHandler = this::sendAccessTokenResponse;
  private AuthenticationFailureHandler authenticationFailureHandler = this::sendErrorResponse;

  /**
   * Constructs an {@code Oauth2ResourceOwnerTokenEndpointFilter} using the provided parameters.
   *
   * @param authenticationManager the authentication manager
   */
  public Oauth2ResourceOwnerTokenEndpointFilter(AuthenticationManager authenticationManager) {
    this(authenticationManager, DEFAULT_TOKEN_ENDPOINT_URI);
  }

  /**
   * Constructs an {@code Oauth2ResourceOwnerTokenEndpointFilter} using the provided parameters.
   *
   * @param authenticationManager the authentication manager
   * @param tokenEndpointUri the endpoint {@code URI} for access token requests
   */
  public Oauth2ResourceOwnerTokenEndpointFilter(
      AuthenticationManager authenticationManager, String tokenEndpointUri) {
    Assert.notNull(authenticationManager, "authenticationManager cannot be null");
    Assert.hasText(tokenEndpointUri, "tokenEndpointUri cannot be empty");
    this.authenticationManager = authenticationManager;
    this.tokenEndpointMatcher = new AntPathRequestMatcher(tokenEndpointUri, HttpMethod.POST.name());
    this.authenticationConverter =
        new DelegatingAuthenticationConverter(
            Arrays.asList(
                new OAuth2AuthorizationCodeAuthenticationConverter(),
                new OAuth2RefreshTokenAuthenticationConverter(),
                new OAuth2ClientCredentialsAuthenticationConverter()));
    OAuth2AccessTokenResponseHttpMessageConverter accessTokenResponseHttpMessageConverter =
        new OAuth2AccessTokenResponseHttpMessageConverter();
    accessTokenResponseHttpMessageConverter.setAccessTokenResponseParametersConverter(
        new RestAccessTokenResponseHttpMessageConverter());
    this.accessTokenHttpResponseConverter = accessTokenResponseHttpMessageConverter;
    OAuth2ErrorHttpMessageConverter errorHttpMessageConverter =
        new OAuth2ErrorHttpMessageConverter();
    errorHttpMessageConverter.setErrorParametersConverter(new RestOAuth2ErrorParametersConverter());
    this.errorHttpResponseConverter = errorHttpMessageConverter;
  }

  private static void throwError(String errorCode, String parameterName) {
    OAuth2Error error =
        new OAuth2Error(errorCode, "OAuth 2.0 Parameter: " + parameterName, DEFAULT_ERROR_URI);
    throw new OAuth2AuthenticationException(error);
  }

  @Override
  protected void doFilterInternal(
      HttpServletRequest request, HttpServletResponse response, FilterChain filterChain)
      throws ServletException, IOException {

    if (!this.tokenEndpointMatcher.matches(request)) {
      filterChain.doFilter(request, response);
      return;
    }

    try {
      Authentication authorizationResourceOwnerAuthentication =
          this.authenticationConverter.convert(request);
      if (authorizationResourceOwnerAuthentication == null) {
        throwError(OAuth2ErrorCodes.UNSUPPORTED_GRANT_TYPE, AuthParameterNames.AUTH_TYPE);
      }
      if (authorizationResourceOwnerAuthentication instanceof AbstractAuthenticationToken) {
        ((AbstractAuthenticationToken) authorizationResourceOwnerAuthentication)
            .setDetails(this.authenticationDetailsSource.buildDetails(request));
      }

      OAuth2AccessTokenAuthenticationToken accessTokenAuthentication =
          (OAuth2AccessTokenAuthenticationToken)
              this.authenticationManager.authenticate(authorizationResourceOwnerAuthentication);
      this.authenticationSuccessHandler.onAuthenticationSuccess(
          request, response, accessTokenAuthentication);
    } catch (OAuth2AuthenticationException ex) {
      SecurityContextHolder.clearContext();
      if (this.logger.isTraceEnabled()) {
        this.logger.trace(LogMessage.format("Token request failed: %s", ex.getError()), ex);
      }
      this.authenticationFailureHandler.onAuthenticationFailure(request, response, ex);
    }
  }

  /**
   * Sets the {@link AuthenticationDetailsSource} used for building an authentication details
   * instance from {@link HttpServletRequest}.
   *
   * @param authenticationDetailsSource the {@link AuthenticationDetailsSource} used for building an
   *     authentication details instance from {@link HttpServletRequest}
   */
  public void setAuthenticationDetailsSource(
      AuthenticationDetailsSource<HttpServletRequest, ?> authenticationDetailsSource) {
    Assert.notNull(authenticationDetailsSource, "authenticationDetailsSource cannot be null");
    this.authenticationDetailsSource = authenticationDetailsSource;
  }

  /**
   * Sets the {@link AuthenticationConverter} used when attempting to extract an Access Token
   * Request from {@link HttpServletRequest} to an instance of {@link
   * OAuth2ResourceOwnerBaseAuthenticationToken} used for authenticating the authorization grant.
   *
   * @param authenticationConverter the {@link AuthenticationConverter} used when attempting to
   *     extract an Access Token Request from {@link HttpServletRequest}
   */
  public void setAuthenticationConverter(AuthenticationConverter authenticationConverter) {
    Assert.notNull(authenticationConverter, "authenticationConverter cannot be null");
    this.authenticationConverter = authenticationConverter;
  }

  /**
   * Sets the {@link AuthenticationSuccessHandler} used for handling an {@link
   * OAuth2AccessTokenAuthenticationToken} and returning the {@link OAuth2AccessTokenResponse Access
   * Token Response}.
   *
   * @param authenticationSuccessHandler the {@link AuthenticationSuccessHandler} used for handling
   *     an {@link OAuth2AccessTokenAuthenticationToken}
   */
  public void setAuthenticationSuccessHandler(
      AuthenticationSuccessHandler authenticationSuccessHandler) {
    Assert.notNull(authenticationSuccessHandler, "authenticationSuccessHandler cannot be null");
    this.authenticationSuccessHandler = authenticationSuccessHandler;
  }

  /**
   * Sets the {@link AuthenticationFailureHandler} used for handling an {@link
   * OAuth2AuthenticationException} and returning the {@link OAuth2Error Error Response}.
   *
   * @param authenticationFailureHandler the {@link AuthenticationFailureHandler} used for handling
   *     an {@link OAuth2AuthenticationException}
   */
  public void setAuthenticationFailureHandler(
      AuthenticationFailureHandler authenticationFailureHandler) {
    Assert.notNull(authenticationFailureHandler, "authenticationFailureHandler cannot be null");
    this.authenticationFailureHandler = authenticationFailureHandler;
  }

  private void sendAccessTokenResponse(
      HttpServletRequest request, HttpServletResponse response, Authentication authentication)
      throws IOException {

    OAuth2AccessTokenAuthenticationToken accessTokenAuthentication =
        (OAuth2AccessTokenAuthenticationToken) authentication;

    OAuth2AccessToken accessToken = accessTokenAuthentication.getAccessToken();
    OAuth2RefreshToken refreshToken = accessTokenAuthentication.getRefreshToken();
    Map<String, Object> additionalParameters = accessTokenAuthentication.getAdditionalParameters();

    OAuth2AccessTokenResponse.Builder builder =
        OAuth2AccessTokenResponse.withToken(accessToken.getTokenValue())
            .tokenType(accessToken.getTokenType())
            .scopes(accessToken.getScopes());
    if (accessToken.getIssuedAt() != null && accessToken.getExpiresAt() != null) {
      builder.expiresIn(
          ChronoUnit.SECONDS.between(accessToken.getIssuedAt(), accessToken.getExpiresAt()));
    }
    if (refreshToken != null) {
      builder.refreshToken(refreshToken.getTokenValue());
    }
    if (!CollectionUtils.isEmpty(additionalParameters)) {
      builder.additionalParameters(additionalParameters);
      String userId = (String) additionalParameters.get(SecurityConstants.USER_ID);
      log.info("user {} login successful.", userId);
      LogUtil.successLogin(userId);
    }

    OAuth2AccessTokenResponse accessTokenResponse = builder.build();
    ServletServerHttpResponse httpResponse = new ServletServerHttpResponse(response);
    this.accessTokenHttpResponseConverter.write(
        accessTokenResponse, MediaType.APPLICATION_JSON, httpResponse);
  }

  private void sendErrorResponse(
      HttpServletRequest request, HttpServletResponse response, AuthenticationException exception)
      throws IOException {

    OAuth2Error error = ((OAuth2AuthenticationException) exception).getError();
    ServletServerHttpResponse httpResponse = new ServletServerHttpResponse(response);
    httpResponse.setStatusCode(HttpStatus.BAD_REQUEST);
    this.errorHttpResponseConverter.write(error, MediaType.APPLICATION_JSON, httpResponse);
  }
}
