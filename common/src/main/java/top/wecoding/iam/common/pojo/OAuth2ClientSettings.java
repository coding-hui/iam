package top.wecoding.iam.common.pojo;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;
import org.springframework.security.oauth2.jose.jws.JwsAlgorithm;
import org.springframework.security.oauth2.jose.jws.MacAlgorithm;
import org.springframework.security.oauth2.jose.jws.SignatureAlgorithm;
import org.springframework.security.oauth2.server.authorization.config.ClientSettings;
import org.springframework.util.StringUtils;

import java.io.Serializable;

/**
 * @author liuyuhui
 * @date 2022/10/5
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class OAuth2ClientSettings implements Serializable {

  private static final long serialVersionUID = -2379117687423147472L;

  private Boolean requireProofKey;

  private Boolean requireAuthorizationConsent;

  private String jwkSetUrl;

  private String signingAlgorithm;

  public static OAuth2ClientSettings fromClientSettings(ClientSettings clientSettings) {
    OAuth2ClientSettings oAuth2ClientSettings = new OAuth2ClientSettings();
    oAuth2ClientSettings.setRequireProofKey(clientSettings.isRequireProofKey());
    oAuth2ClientSettings.setRequireAuthorizationConsent(
        clientSettings.isRequireAuthorizationConsent());
    oAuth2ClientSettings.setJwkSetUrl(clientSettings.getJwkSetUrl());
    JwsAlgorithm algorithm = clientSettings.getTokenEndpointAuthenticationSigningAlgorithm();
    if (algorithm != null) {
      oAuth2ClientSettings.setSigningAlgorithm(algorithm.getName());
    }
    return oAuth2ClientSettings;
  }

  public ClientSettings toClientSettings() {
    ClientSettings.Builder builder =
        ClientSettings.builder()
            .requireProofKey(this.requireProofKey)
            .requireAuthorizationConsent(this.requireAuthorizationConsent);
    SignatureAlgorithm signatureAlgorithm = SignatureAlgorithm.from(this.signingAlgorithm);
    JwsAlgorithm jwsAlgorithm =
        signatureAlgorithm == null ? MacAlgorithm.from(this.signingAlgorithm) : signatureAlgorithm;
    if (jwsAlgorithm != null) {
      builder.tokenEndpointAuthenticationSigningAlgorithm(jwsAlgorithm);
    }
    if (StringUtils.hasText(this.jwkSetUrl)) {
      builder.jwkSetUrl(jwkSetUrl);
    }
    return builder.build();
  }
}
