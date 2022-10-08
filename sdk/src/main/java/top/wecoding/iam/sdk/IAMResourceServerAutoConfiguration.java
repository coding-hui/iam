package top.wecoding.iam.sdk;

import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.security.oauth2.server.resource.introspection.OpaqueTokenIntrospector;
import top.wecoding.iam.sdk.introspection.IAMOpaqueTokenIntrospector;
import top.wecoding.iam.sdk.web.IAMBearerTokenExtractor;
import top.wecoding.iam.sdk.web.PermissionService;
import top.wecoding.iam.sdk.web.ResourceAuthExceptionEntryPoint;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@RequiredArgsConstructor
@EnableConfigurationProperties(IgnoreWhiteProperties.class)
public class IAMResourceServerAutoConfiguration {

  @Bean("pms")
  public PermissionService permissionService() {
    return new PermissionService();
  }

  @Bean
  public IAMBearerTokenExtractor iamBearerTokenExtractor(IgnoreWhiteProperties urlProperties) {
    return new IAMBearerTokenExtractor(urlProperties);
  }

  @Bean
  public ResourceAuthExceptionEntryPoint resourceAuthExceptionEntryPoint(
      ObjectMapper objectMapper) {
    return new ResourceAuthExceptionEntryPoint(objectMapper);
  }

  @Bean
  public OpaqueTokenIntrospector opaqueTokenIntrospector(
      OAuth2AuthorizationService authorizationService) {
    return new IAMOpaqueTokenIntrospector(authorizationService);
  }
}
