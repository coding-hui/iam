package top.wecoding.iam.server.security.authorization.authentication;

import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.context.MessageSource;
import org.springframework.context.support.MessageSourceAccessor;
import org.springframework.context.support.ReloadableResourceBundleMessageSource;
import org.springframework.security.authentication.*;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.oauth2.core.*;
import org.springframework.security.oauth2.server.authorization.OAuth2Authorization;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.security.oauth2.server.authorization.OAuth2TokenType;
import org.springframework.security.oauth2.server.authorization.authentication.OAuth2AccessTokenAuthenticationToken;
import org.springframework.security.oauth2.server.authorization.authentication.OAuth2ClientAuthenticationToken;
import org.springframework.security.oauth2.server.authorization.client.RegisteredClient;
import org.springframework.security.oauth2.server.authorization.context.AuthorizationServerContextHolder;
import org.springframework.security.oauth2.server.authorization.token.DefaultOAuth2TokenContext;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenContext;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenGenerator;
import org.springframework.util.Assert;
import top.wecoding.iam.common.constant.OAuth2ErrorCodesExpand;
import top.wecoding.iam.common.enums.IamErrorCode;
import top.wecoding.iam.common.model.request.LoginRequest;
import top.wecoding.iam.framework.cache.UserFailCountCacheKeyBuilder;
import top.wecoding.iam.framework.exception.ScopeException;
import top.wecoding.iam.framework.props.AppProperties;
import top.wecoding.redis.util.RedisUtils;

import java.security.Principal;
import java.time.Duration;
import java.time.Instant;
import java.util.Locale;
import java.util.Objects;
import java.util.Set;
import java.util.function.Supplier;

/**
 * 处理自定义授权
 *
 * @author liuyuhui
 * @date 2022/10/3
 */
public abstract class OAuth2ResourceOwnerBaseAuthenticationProvider<
        T extends OAuth2ResourceOwnerBaseAuthenticationToken>
    implements AuthenticationProvider {

  private static final Logger LOGGER =
      LogManager.getLogger(OAuth2ResourceOwnerBaseAuthenticationProvider.class);

  private static final String ERROR_URI =
      "https://datatracker.ietf.org/doc/html/rfc6749#section-4.1.2.1";

  private final OAuth2AuthorizationService authorizationService;

  private final OAuth2TokenGenerator<? extends OAuth2Token> tokenGenerator;

  private final AuthenticationManager authenticationManager;

  private final MessageSourceAccessor messages;

  @Deprecated private Supplier<String> refreshTokenGenerator;

  /**
   * Constructs an {@code OAuth2AuthorizationCodeAuthenticationProvider} using the provided
   * parameters.
   *
   * @param authorizationService the authorization service
   * @param tokenGenerator the token generator
   * @since 0.2.3
   */
  public OAuth2ResourceOwnerBaseAuthenticationProvider(
      AuthenticationManager authenticationManager,
      OAuth2AuthorizationService authorizationService,
      OAuth2TokenGenerator<? extends OAuth2Token> tokenGenerator) {
    Assert.notNull(authorizationService, "authorizationService cannot be null");
    Assert.notNull(tokenGenerator, "tokenGenerator cannot be null");
    this.authenticationManager = authenticationManager;
    this.authorizationService = authorizationService;
    this.tokenGenerator = tokenGenerator;
    this.messages = new MessageSourceAccessor(iamMessageSource(), Locale.CHINA);
  }

  @Deprecated
  public void setRefreshTokenGenerator(Supplier<String> refreshTokenGenerator) {
    Assert.notNull(refreshTokenGenerator, "refreshTokenGenerator cannot be null");
    this.refreshTokenGenerator = refreshTokenGenerator;
  }

  public abstract UsernamePasswordAuthenticationToken buildToken(LoginRequest loginRequest);

  @Override
  public abstract boolean supports(Class<?> authentication);

  public abstract void checkClient(RegisteredClient registeredClient);

  @Override
  public Authentication authenticate(Authentication authentication) throws AuthenticationException {

    var resourceOwnerBaseAuthentication = (T) authentication;

    OAuth2ClientAuthenticationToken clientPrincipal =
        getAuthenticatedClientElseThrowInvalidClient(resourceOwnerBaseAuthentication);

    RegisteredClient registeredClient = clientPrincipal.getRegisteredClient();
    checkClient(registeredClient);

    Set<String> authorizedScopes = resourceOwnerBaseAuthentication.getScopes();

    UsernamePasswordAuthenticationToken usernamePasswordAuthenticationToken =
        buildToken(resourceOwnerBaseAuthentication.getLoginRequest());

    try {

      LOGGER.debug(
          "got usernamePasswordAuthenticationToken=" + usernamePasswordAuthenticationToken);

      Authentication usernamePasswordAuthentication =
          authenticationManager.authenticate(usernamePasswordAuthenticationToken);

      incrFailCounter(authentication, usernamePasswordAuthenticationToken);

      DefaultOAuth2TokenContext.Builder tokenContextBuilder =
          DefaultOAuth2TokenContext.builder()
              .registeredClient(registeredClient)
              .principal(usernamePasswordAuthentication)
              .authorizationServerContext(AuthorizationServerContextHolder.getContext())
              .authorizedScopes(authorizedScopes)
              .authorizationGrantType(AuthorizationGrantType.PASSWORD)
              .authorizationGrant(resourceOwnerBaseAuthentication);

      OAuth2Authorization.Builder authorizationBuilder =
          OAuth2Authorization.withRegisteredClient(registeredClient)
              .principalName(usernamePasswordAuthentication.getName())
              .authorizationGrantType(AuthorizationGrantType.PASSWORD)
              .authorizedScopes(authorizedScopes);

      // ----- Access token -----
      OAuth2TokenContext tokenContext =
          tokenContextBuilder.tokenType(OAuth2TokenType.ACCESS_TOKEN).build();
      OAuth2Token generatedAccessToken = this.tokenGenerator.generate(tokenContext);
      if (generatedAccessToken == null) {
        OAuth2Error error =
            new OAuth2Error(
                OAuth2ErrorCodes.SERVER_ERROR,
                "The token generator failed to generate the access token.",
                ERROR_URI);
        throw new OAuth2AuthenticationException(error);
      }
      OAuth2AccessToken accessToken =
          new OAuth2AccessToken(
              OAuth2AccessToken.TokenType.BEARER,
              generatedAccessToken.getTokenValue(),
              generatedAccessToken.getIssuedAt(),
              generatedAccessToken.getExpiresAt(),
              tokenContext.getAuthorizedScopes());
      if (generatedAccessToken instanceof ClaimAccessor) {
        authorizationBuilder
            .id(accessToken.getTokenValue())
            .token(
                accessToken,
                (metadata) ->
                    metadata.put(
                        OAuth2Authorization.Token.CLAIMS_METADATA_NAME,
                        ((ClaimAccessor) generatedAccessToken).getClaims()))
            .authorizedScopes(authorizedScopes)
            .attribute(Principal.class.getName(), usernamePasswordAuthentication);
      } else {
        authorizationBuilder.id(accessToken.getTokenValue()).accessToken(accessToken);
      }

      // ----- Refresh token -----
      OAuth2RefreshToken refreshToken = null;
      if (registeredClient
              .getAuthorizationGrantTypes()
              .contains(AuthorizationGrantType.REFRESH_TOKEN)
          &&
          // Do not issue refresh token to public client
          !clientPrincipal
              .getClientAuthenticationMethod()
              .equals(ClientAuthenticationMethod.NONE)) {

        if (this.refreshTokenGenerator != null) {
          Instant issuedAt = Instant.now();
          Instant expiresAt =
              issuedAt.plus(registeredClient.getTokenSettings().getRefreshTokenTimeToLive());
          refreshToken =
              new OAuth2RefreshToken(this.refreshTokenGenerator.get(), issuedAt, expiresAt);
        } else {
          tokenContext = tokenContextBuilder.tokenType(OAuth2TokenType.REFRESH_TOKEN).build();
          OAuth2Token generatedRefreshToken = this.tokenGenerator.generate(tokenContext);
          if (!(generatedRefreshToken instanceof OAuth2RefreshToken)) {
            OAuth2Error error =
                new OAuth2Error(
                    OAuth2ErrorCodes.SERVER_ERROR,
                    "The token generator failed to generate the refresh token.",
                    ERROR_URI);
            throw new OAuth2AuthenticationException(error);
          }
          refreshToken = (OAuth2RefreshToken) generatedRefreshToken;
        }
        authorizationBuilder.refreshToken(refreshToken);
      }

      OAuth2Authorization authorization = authorizationBuilder.build();

      this.authorizationService.save(authorization);

      LOGGER.debug("returning OAuth2AccessTokenAuthenticationToken");

      refreshFailCount(authentication, usernamePasswordAuthenticationToken);

      return new OAuth2AccessTokenAuthenticationToken(
          registeredClient,
          clientPrincipal,
          accessToken,
          refreshToken,
          Objects.requireNonNull(authorization.getAccessToken().getClaims()));

    } catch (Exception ex) {
      LOGGER.error("problem in authenticate", ex);

      incrFailCounter(authentication, usernamePasswordAuthenticationToken);

      throw oAuth2AuthenticationException(authentication, (AuthenticationException) ex);
    }
  }

  private OAuth2AuthenticationException oAuth2AuthenticationException(
      Authentication authentication, AuthenticationException authenticationException) {
    if (authenticationException instanceof UsernameNotFoundException) {
      return new OAuth2AuthenticationException(
          new OAuth2Error(
              OAuth2ErrorCodesExpand.USERNAME_NOT_FOUND,
              this.messages.getMessage(
                  "JdbcDaoImpl.notFound",
                  new Object[] {authentication.getName()},
                  "Username {0} not found"),
              ""));
    }
    if (authenticationException instanceof BadCredentialsException) {
      return new OAuth2AuthenticationException(
          new OAuth2Error(
              OAuth2ErrorCodesExpand.BAD_CREDENTIALS,
              this.messages.getMessage(
                  "AbstractUserDetailsAuthenticationProvider.badCredentials", "Bad credentials"),
              ""));
    }
    if (authenticationException instanceof LockedException) {
      return new OAuth2AuthenticationException(
          new OAuth2Error(
              OAuth2ErrorCodesExpand.USER_LOCKED,
              this.messages.getMessage(
                  "AbstractUserDetailsAuthenticationProvider.locked", "User account is locked"),
              ""));
    }
    if (authenticationException instanceof DisabledException) {
      return new OAuth2AuthenticationException(
          new OAuth2Error(
              OAuth2ErrorCodesExpand.USER_DISABLE,
              this.messages.getMessage(
                  "AbstractUserDetailsAuthenticationProvider.disabled", "User is disabled"),
              ""));
    }
    if (authenticationException instanceof AccountExpiredException) {
      return new OAuth2AuthenticationException(
          new OAuth2Error(
              OAuth2ErrorCodesExpand.USER_EXPIRED,
              this.messages.getMessage(
                  "AbstractUserDetailsAuthenticationProvider.expired", "User account has expired"),
              ""));
    }
    if (authenticationException instanceof CredentialsExpiredException) {
      return new OAuth2AuthenticationException(
          new OAuth2Error(
              OAuth2ErrorCodesExpand.CREDENTIALS_EXPIRED,
              this.messages.getMessage(
                  "AbstractUserDetailsAuthenticationProvider.credentialsExpired",
                  "User credentials have expired"),
              ""));
    }
    if (authenticationException instanceof ScopeException) {
      return new OAuth2AuthenticationException(
          new OAuth2Error(
              OAuth2ErrorCodes.INVALID_SCOPE,
              this.messages.getMessage(
                  "AbstractAccessDecisionManager.accessDenied", "invalid_scope"),
              ""));
    }
    return new OAuth2AuthenticationException(
        new OAuth2Error(
            OAuth2ErrorCodes.INVALID_REQUEST,
            this.messages.getMessage(IamErrorCode.UNAUTHORIZED.getCode()),
            ""));
  }

  private OAuth2ClientAuthenticationToken getAuthenticatedClientElseThrowInvalidClient(
      Authentication authentication) {

    OAuth2ClientAuthenticationToken clientPrincipal = null;

    if (OAuth2ClientAuthenticationToken.class.isAssignableFrom(
        authentication.getPrincipal().getClass())) {
      clientPrincipal = (OAuth2ClientAuthenticationToken) authentication.getPrincipal();
    }

    if (clientPrincipal != null && clientPrincipal.isAuthenticated()) {
      return clientPrincipal;
    }

    throw new OAuth2AuthenticationException(OAuth2ErrorCodes.INVALID_CLIENT);
  }

  private void incrFailCounter(
      Authentication clientAuthentication,
      UsernamePasswordAuthenticationToken usernamePasswordAuthenticationToken) {
    String clientId = clientAuthentication.getName();
    String username = usernamePasswordAuthenticationToken.getName();
    String key = new UserFailCountCacheKeyBuilder().build(clientId, username).getKey();

    Long num = RedisUtils.incr(key);

    if (Objects.nonNull(num) && num > AppProperties.getUserFailCount()) {
      Long ttl = RedisUtils.ttl(key);
      throw new OAuth2AuthenticationException(
          new OAuth2Error(
              IamErrorCode.TOO_MANY_FAILURES_PLEASE_TRY_AGAIN_LATER.getCode(),
              this.messages.getMessage(
                  IamErrorCode.TOO_MANY_FAILURES_PLEASE_TRY_AGAIN_LATER.getCode(),
                  new Object[] {ttl / 60}),
              ""));
    }
    RedisUtils.expire(key, Duration.ofSeconds(AppProperties.getUserFailLockTime()));
  }

  private void refreshFailCount(
      Authentication clientAuthentication,
      UsernamePasswordAuthenticationToken usernamePasswordAuthenticationToken) {
    String clientId = clientAuthentication.getName();
    String username = usernamePasswordAuthenticationToken.getName();
    String key = new UserFailCountCacheKeyBuilder().build(clientId, username).getKey();

    RedisUtils.del(key);
  }

  private MessageSource iamMessageSource() {
    ReloadableResourceBundleMessageSource messageSource =
        new ReloadableResourceBundleMessageSource();
    messageSource.addBasenames("classpath:i18n/errors/messages");
    messageSource.addBasenames("classpath:i18n/errors/messages-iam");
    messageSource.addBasenames("classpath:i18n/errors/messages-common");
    messageSource.setDefaultLocale(Locale.CHINA);
    return messageSource;
  }
}
