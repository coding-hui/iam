package top.wecoding.iam.sdk.feign;

import feign.RequestInterceptor;
import org.springframework.context.annotation.Bean;
import org.springframework.security.oauth2.server.resource.web.BearerTokenResolver;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
public class IAMFeignClientConfiguration {

  @Bean
  public RequestInterceptor oauthRequestInterceptor(BearerTokenResolver tokenResolver) {
    return new IAMFeignRequestInterceptor(tokenResolver);
  }
}
