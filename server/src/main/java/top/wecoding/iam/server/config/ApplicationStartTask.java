package top.wecoding.iam.server.config;

import com.baomidou.mybatisplus.core.toolkit.IdWorker;
import java.time.Duration;
import java.util.Arrays;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.ApplicationRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.oauth2.core.AuthorizationGrantType;
import org.springframework.security.oauth2.core.ClientAuthenticationMethod;
import org.springframework.security.oauth2.core.OAuth2TokenFormat;
import org.springframework.security.oauth2.jose.jws.SignatureAlgorithm;
import org.springframework.util.StringUtils;
import top.wecoding.core.enums.iam.UserTypeEnum;
import top.wecoding.iam.common.constant.SecurityConstants;
import top.wecoding.iam.common.entity.OAuth2ClientSettings;
import top.wecoding.iam.common.entity.OAuth2TokenSettings;
import top.wecoding.iam.server.entity.Oauth2Client;
import top.wecoding.iam.server.entity.User;
import top.wecoding.iam.server.enums.UserStateEnum;
import top.wecoding.iam.server.mapper.Oauth2ClientMapper;
import top.wecoding.iam.server.mapper.UserMapper;
import top.wecoding.iam.server.util.PasswordEncoderUtil;

/**
 * @author Liuyuhui
 * @date 2022/10/8
 */
@Slf4j
@Configuration
@RequiredArgsConstructor
public class ApplicationStartTask {

  private final UserMapper userMapper;

  private final Oauth2ClientMapper oauth2ClientMapper;

  @Bean
  public ApplicationRunner initialize() {
    return args -> {
      initializeDefaultClient();
      initializeDefaultSuperAdministrator();
    };
  }

  private void initializeDefaultClient() {
    String clientId = "wecoding";
    Oauth2Client client = oauth2ClientMapper.getByClientId(clientId);
    if (null != client) {
      return;
    }
    client =
        Oauth2Client.builder()
            .clientId(clientId)
            .clientName(clientId)
            .clientSecret(clientId)
            .clientAuthenticationMethods(
                StringUtils.collectionToCommaDelimitedString(
                    Arrays.asList(
                        ClientAuthenticationMethod.CLIENT_SECRET_POST.getValue(),
                        ClientAuthenticationMethod.CLIENT_SECRET_BASIC.getValue(),
                        ClientAuthenticationMethod.CLIENT_SECRET_JWT.getValue())))
            .authorizationGrantTypes(
                StringUtils.collectionToCommaDelimitedString(
                    Arrays.asList(
                        AuthorizationGrantType.AUTHORIZATION_CODE.getValue(),
                        AuthorizationGrantType.CLIENT_CREDENTIALS.getValue(),
                        AuthorizationGrantType.REFRESH_TOKEN.getValue(),
                        AuthorizationGrantType.JWT_BEARER.getValue(),
                        AuthorizationGrantType.PASSWORD.getValue())))
            .redirectUris(SecurityConstants.PROJECT_LICENSE)
            .scopes("server")
            .clientSettings(
                OAuth2ClientSettings.builder()
                    .requireProofKey(true)
                    .requireAuthorizationConsent(false)
                    .signingAlgorithm(SignatureAlgorithm.RS256.getName())
                    .build())
            .tokenSettings(
                OAuth2TokenSettings.builder()
                    .accessTokenTimeToLive(Duration.ofHours(6).getSeconds())
                    .refreshTokenTimeToLive(Duration.ofDays(30).getSeconds())
                    .tokenFormat(OAuth2TokenFormat.REFERENCE.getValue())
                    .reuseRefreshTokens(true)
                    .build())
            .build();
    oauth2ClientMapper.insert(client);
    log.info("initialize default client done");
  }

  private void initializeDefaultSuperAdministrator() {
    String tenantId = String.valueOf(Long.MAX_VALUE);
    String username = "ADMIN";
    User user = userMapper.getByTenantIdAndUsername(tenantId, username);
    if (null != user) {
      return;
    }
    user =
        User.builder()
            .userId(IdWorker.getIdStr())
            .tenantId(tenantId)
            .username(username)
            .nickName(username)
            .country("")
            .email("wecoding@yeah.net")
            .password(PasswordEncoderUtil.encode("WECODING"))
            .userType(UserTypeEnum.LOCAL.code())
            .userState(UserStateEnum.DEFAULT.code())
            .defaultPwd(true)
            .build();
    userMapper.insert(user);
    log.info("initialize default super administrator done");
  }
}
