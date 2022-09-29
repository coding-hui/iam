package top.wecoding.iam.sdk.config;

import lombok.AllArgsConstructor;
import org.springframework.boot.autoconfigure.AutoConfiguration;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.boot.web.servlet.FilterRegistrationBean;
import org.springframework.context.annotation.Bean;
import org.springframework.core.annotation.Order;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.lang.NonNull;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
import top.wecoding.iam.sdk.aspect.InnerAuthAspect;
import top.wecoding.iam.sdk.filter.AuthenticationProcessingFilter;
import top.wecoding.iam.sdk.interceptor.ClientInterceptor;
import top.wecoding.iam.sdk.interceptor.SecurityInterceptor;
import top.wecoding.iam.sdk.props.SecurityProperties;
import top.wecoding.iam.common.provider.ClientDetailsService;
import top.wecoding.iam.common.provider.client.JdbcClientDetailsService;
import top.wecoding.iam.sdk.registry.SecurityRegistry;

/**
 * 安全配置类
 *
 * @author liuyuhui
 * @date 2022
 * @qq 1515418211
 */
@Order
@AutoConfiguration
@AllArgsConstructor
@EnableConfigurationProperties({SecurityProperties.class})
public class SecurityConfiguration implements WebMvcConfigurer {

  private final JdbcTemplate jdbcTemplate;
  private final SecurityRegistry securityRegistry;
  private final SecurityProperties securityProperties;

  @Override
  public void addInterceptors(@NonNull InterceptorRegistry registry) {
    securityProperties
        .getClient()
        .forEach(
            cs ->
                registry
                    .addInterceptor(new ClientInterceptor(cs.getClientId()))
                    .addPathPatterns(cs.getPathPatterns()));

    if (securityRegistry.isEnabled()) {
      registry
          .addInterceptor(new SecurityInterceptor())
          .addPathPatterns(securityRegistry.getIncludePatterns())
          .excludePathPatterns(securityRegistry.getExcludePatterns())
          .excludePathPatterns(securityRegistry.getDefaultExcludePatterns())
          .excludePathPatterns(securityProperties.getWhites());
    }
  }

  @Bean
  public FilterRegistrationBean<AuthenticationProcessingFilter>
      authenticationProcessingFilterFilterRegistrationBean() {
    FilterRegistrationBean<AuthenticationProcessingFilter> registrationBean =
        new FilterRegistrationBean<>();
    registrationBean.setFilter(new AuthenticationProcessingFilter());
    registrationBean.setUrlPatterns(securityRegistry.getIncludePatterns());
    registrationBean.setOrder(securityRegistry.getOrder());
    return registrationBean;
  }

  @Bean
  @ConditionalOnMissingBean(ClientDetailsService.class)
  public ClientDetailsService clientDetailsService() {
    return new JdbcClientDetailsService(jdbcTemplate);
  }

  @Bean
  public InnerAuthAspect innerAuthAspect() {
    return new InnerAuthAspect();
  }
}
