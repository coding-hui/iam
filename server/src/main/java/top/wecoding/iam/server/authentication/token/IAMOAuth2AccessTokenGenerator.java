package top.wecoding.iam.server.authentication.token;

import org.springframework.lang.Nullable;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.crypto.keygen.Base64StringKeyGenerator;
import org.springframework.security.crypto.keygen.StringKeyGenerator;
import org.springframework.security.oauth2.core.ClaimAccessor;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.security.oauth2.core.OAuth2TokenFormat;
import org.springframework.security.oauth2.core.OAuth2TokenType;
import org.springframework.security.oauth2.core.endpoint.OAuth2ParameterNames;
import org.springframework.security.oauth2.server.authorization.client.RegisteredClient;
import org.springframework.security.oauth2.server.authorization.token.*;
import org.springframework.util.Assert;
import org.springframework.util.CollectionUtils;
import org.springframework.util.StringUtils;
import top.wecoding.iam.common.constant.UserConstant;

import java.time.Instant;
import java.util.*;
import java.util.concurrent.ThreadLocalRandom;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
public class IAMOAuth2AccessTokenGenerator implements OAuth2TokenGenerator<OAuth2AccessToken> {

  private final StringKeyGenerator accessTokenGenerator =
      new Base64StringKeyGenerator(Base64.getUrlEncoder().withoutPadding(), 96);

  private OAuth2TokenCustomizer<OAuth2TokenClaimsContext> accessTokenCustomizer;

  @Nullable
  @Override
  public OAuth2AccessToken generate(OAuth2TokenContext context) {
    if (!OAuth2TokenType.ACCESS_TOKEN.equals(context.getTokenType())
        || !OAuth2TokenFormat.REFERENCE.equals(
            context.getRegisteredClient().getTokenSettings().getAccessTokenFormat())) {
      return null;
    }

    String issuer = null;
    if (context.getProviderContext() != null) {
      issuer = context.getProviderContext().getIssuer();
    }
    RegisteredClient registeredClient = context.getRegisteredClient();

    Instant issuedAt = Instant.now();
    Instant expiresAt =
        issuedAt.plus(registeredClient.getTokenSettings().getAccessTokenTimeToLive());

    // @formatter:off
    OAuth2TokenClaimsSet.Builder claimsBuilder = OAuth2TokenClaimsSet.builder();
    if (StringUtils.hasText(issuer)) {
      claimsBuilder.issuer(issuer);
    }
    claimsBuilder
        .subject(context.getPrincipal().getName())
        .audience(Collections.singletonList(registeredClient.getClientId()))
        .issuedAt(issuedAt)
        .expiresAt(expiresAt)
        .notBefore(issuedAt)
        .id(UUID.randomUUID().toString());
    if (!CollectionUtils.isEmpty(context.getAuthorizedScopes())) {
      claimsBuilder.claim(OAuth2ParameterNames.SCOPE, context.getAuthorizedScopes());
    }
    // @formatter:on

    if (this.accessTokenCustomizer != null) {
      // @formatter:off
      OAuth2TokenClaimsContext.Builder accessTokenContextBuilder =
          OAuth2TokenClaimsContext.with(claimsBuilder)
              .registeredClient(context.getRegisteredClient())
              .principal(context.getPrincipal())
              .providerContext(context.getProviderContext())
              .authorizedScopes(context.getAuthorizedScopes())
              .tokenType(context.getTokenType())
              .authorizationGrantType(context.getAuthorizationGrantType());
      if (context.getAuthorization() != null) {
        accessTokenContextBuilder.authorization(context.getAuthorization());
      }
      if (context.getAuthorizationGrant() != null) {
        accessTokenContextBuilder.authorizationGrant(context.getAuthorizationGrant());
      }
      // @formatter:on

      OAuth2TokenClaimsContext accessTokenContext = accessTokenContextBuilder.build();
      this.accessTokenCustomizer.customize(accessTokenContext);
    }

    OAuth2TokenClaimsSet accessTokenClaimsSet = claimsBuilder.build();

    // key client::username::uuid
    ThreadLocalRandom random = ThreadLocalRandom.current();
    String uuid =
        new UUID(random.nextLong(), random.nextLong())
            .toString()
            .replace(UserConstant.DASH, UserConstant.EMPTY);
    String key =
        String.format(
            "%s:%s:%s",
            SecurityContextHolder.getContext().getAuthentication().getPrincipal(),
            context.getPrincipal().getName(),
            uuid);

    return new IAMOAuth2AccessTokenGenerator.OAuth2AccessTokenClaims(
        OAuth2AccessToken.TokenType.BEARER,
        key,
        accessTokenClaimsSet.getIssuedAt(),
        accessTokenClaimsSet.getExpiresAt(),
        context.getAuthorizedScopes(),
        accessTokenClaimsSet.getClaims());
  }

  /**
   * Sets the {@link OAuth2TokenCustomizer} that customizes the {@link
   * OAuth2TokenClaimsContext#getClaims() claims} for the {@link OAuth2AccessToken}.
   *
   * @param accessTokenCustomizer the {@link OAuth2TokenCustomizer} that customizes the claims for
   *     the {@code OAuth2AccessToken}
   */
  public void setAccessTokenCustomizer(
      OAuth2TokenCustomizer<OAuth2TokenClaimsContext> accessTokenCustomizer) {
    Assert.notNull(accessTokenCustomizer, "accessTokenCustomizer cannot be null");
    this.accessTokenCustomizer = accessTokenCustomizer;
  }

  private static final class OAuth2AccessTokenClaims extends OAuth2AccessToken
      implements ClaimAccessor {

    private final Map<String, Object> claims;

    private OAuth2AccessTokenClaims(
        TokenType tokenType,
        String tokenValue,
        Instant issuedAt,
        Instant expiresAt,
        Set<String> scopes,
        Map<String, Object> claims) {
      super(tokenType, tokenValue, issuedAt, expiresAt, scopes);
      this.claims = claims;
    }

    @Override
    public Map<String, Object> getClaims() {
      return this.claims;
    }
  }
}
