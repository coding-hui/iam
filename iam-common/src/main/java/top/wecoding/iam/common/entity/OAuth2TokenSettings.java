package top.wecoding.iam.common.entity;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.io.Serial;
import java.io.Serializable;
import java.time.Duration;
import java.util.Optional;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.security.oauth2.jose.jws.SignatureAlgorithm;
import org.springframework.security.oauth2.server.authorization.settings.OAuth2TokenFormat;
import org.springframework.security.oauth2.server.authorization.settings.TokenSettings;

/**
 * @author liuyuhui
 * @date 2022/10/5
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class OAuth2TokenSettings implements Serializable {

  @Serial private static final long serialVersionUID = -4884699881799394535L;

  private static final Duration DEFAULT_ACCESS_TOKEN_TIME = Duration.ofHours(2L);

  private static final Duration DEFAULT_REFRESH_TOKEN_TIME = Duration.ofDays(30);

  @JsonProperty("access_token_time_to_live")
  private Long accessTokenTimeToLive;

  @JsonProperty("token_format")
  private String tokenFormat;

  @JsonProperty("reuse_refresh_tokens")
  private Boolean reuseRefreshTokens = true;

  @JsonProperty("refresh_token_time_to_live")
  private Long refreshTokenTimeToLive;

  @JsonProperty("id_token_signature_algorithm")
  private String idTokenSignatureAlgorithm;

  public static OAuth2TokenSettings fromTokenSettings(TokenSettings tokenSettings) {
    OAuth2TokenSettings oAuth2TokenSettings = new OAuth2TokenSettings();
    oAuth2TokenSettings.setAccessTokenTimeToLive(
        tokenSettings.getAccessTokenTimeToLive().getSeconds());
    oAuth2TokenSettings.setTokenFormat(tokenSettings.getAccessTokenFormat().getValue());
    oAuth2TokenSettings.setReuseRefreshTokens(tokenSettings.isReuseRefreshTokens());
    oAuth2TokenSettings.setRefreshTokenTimeToLive(
        tokenSettings.getRefreshTokenTimeToLive().getSeconds());
    oAuth2TokenSettings.setIdTokenSignatureAlgorithm(
        tokenSettings.getIdTokenSignatureAlgorithm().getName());
    return oAuth2TokenSettings;
  }

  public TokenSettings toTokenSettings() {
    return TokenSettings.builder()
        .accessTokenTimeToLive(
            Optional.ofNullable(this.accessTokenTimeToLive)
                .map(Duration::ofSeconds)
                .orElse(DEFAULT_ACCESS_TOKEN_TIME))
        .accessTokenFormat(
            Optional.ofNullable(tokenFormat)
                .map(OAuth2TokenFormat::new)
                .orElse(OAuth2TokenFormat.SELF_CONTAINED))
        .reuseRefreshTokens(this.reuseRefreshTokens)
        .refreshTokenTimeToLive(
            Optional.ofNullable(this.refreshTokenTimeToLive)
                .map(Duration::ofSeconds)
                .orElse(DEFAULT_REFRESH_TOKEN_TIME))
        .idTokenSignatureAlgorithm(
            Optional.ofNullable(idTokenSignatureAlgorithm)
                .map(SignatureAlgorithm::from)
                .orElse(SignatureAlgorithm.RS256))
        .build();
  }
}
