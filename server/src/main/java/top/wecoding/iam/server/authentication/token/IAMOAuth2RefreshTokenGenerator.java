package top.wecoding.iam.server.authentication.token;

import java.time.Instant;
import java.util.UUID;
import java.util.concurrent.ThreadLocalRandom;
import org.springframework.lang.Nullable;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.core.OAuth2RefreshToken;
import org.springframework.security.oauth2.core.OAuth2TokenType;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenContext;
import org.springframework.security.oauth2.server.authorization.token.OAuth2TokenGenerator;
import top.wecoding.iam.common.constant.UserConstant;

/**
 * @author liuyuhui
 * @date 2022/10/4
 * @qq 1515418211
 */
public class IAMOAuth2RefreshTokenGenerator implements OAuth2TokenGenerator<OAuth2RefreshToken> {

  @Nullable
  @Override
  public OAuth2RefreshToken generate(OAuth2TokenContext context) {
    if (!OAuth2TokenType.REFRESH_TOKEN.equals(context.getTokenType())) {
      return null;
    }
    Instant issuedAt = Instant.now();
    Instant expiresAt =
        issuedAt.plus(context.getRegisteredClient().getTokenSettings().getRefreshTokenTimeToLive());

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

    return new OAuth2RefreshToken(key, issuedAt, expiresAt);
  }
}
