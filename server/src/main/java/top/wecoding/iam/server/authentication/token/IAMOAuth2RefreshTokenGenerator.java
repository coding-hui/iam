package top.wecoding.iam.server.authentication.token;

import org.springframework.lang.Nullable;
import org.springframework.security.crypto.keygen.Base64StringKeyGenerator;
import org.springframework.security.crypto.keygen.StringKeyGenerator;
import org.springframework.security.oauth2.core.OAuth2RefreshToken;
import org.springframework.security.oauth2.core.OAuth2TokenType;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenContext;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenGenerator;

import java.time.Instant;
import java.util.Base64;

/**
 * @author liuyuhui
 * @date 2022/10/4
 * @qq 1515418211
 */
public class IAMOAuth2RefreshTokenGenerator implements OAuth2TokenGenerator<OAuth2RefreshToken> {

  private final StringKeyGenerator refreshTokenGenerator =
      new Base64StringKeyGenerator(Base64.getUrlEncoder().withoutPadding(), 32);

  @Nullable
  @Override
  public OAuth2RefreshToken generate(OAuth2TokenContext context) {
    if (!OAuth2TokenType.REFRESH_TOKEN.equals(context.getTokenType())) {
      return null;
    }
    Instant issuedAt = Instant.now();
    Instant expiresAt =
        issuedAt.plus(context.getRegisteredClient().getTokenSettings().getRefreshTokenTimeToLive());
    return new OAuth2RefreshToken(this.refreshTokenGenerator.generateKey(), issuedAt, expiresAt);
  }
}
