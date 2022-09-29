package top.wecoding.iam.sdk.config;

import org.springframework.boot.autoconfigure.AutoConfiguration;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.context.annotation.Bean;
import org.springframework.core.annotation.Order;
import top.wecoding.iam.sdk.registry.SecurityRegistry;

/**
 * @author liuyuhui
 * @date 2022/6/6
 * @qq 1515418211
 */
@Order
@AutoConfiguration(before = SecurityConfiguration.class)
public class RegistryConfiguration {

  /** 服务可以自定义放行接口 */
  @Bean
  @ConditionalOnMissingBean(SecurityRegistry.class)
  public SecurityRegistry securityRegistry() {
    return new SecurityRegistry();
  }
}
