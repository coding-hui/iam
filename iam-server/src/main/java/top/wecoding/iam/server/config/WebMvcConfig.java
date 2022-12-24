package top.wecoding.iam.server.config;

import jakarta.annotation.PostConstruct;
import java.util.List;
import java.util.Locale;
import lombok.RequiredArgsConstructor;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.MessageSource;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.support.ReloadableResourceBundleMessageSource;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.data.redis.serializer.RedisSerializer;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.web.method.support.HandlerMethodArgumentResolver;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
import top.wecoding.iam.framework.props.AppProperties;
import top.wecoding.iam.server.web.UnderlineToCamelArgumentResolver;
import top.wecoding.redis.util.RedisUtils;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
@EnableWebSecurity
@RequiredArgsConstructor
@Configuration(proxyBeanMethods = false)
@EnableConfigurationProperties(AppProperties.class)
public class WebMvcConfig implements WebMvcConfigurer {

  private final RedisTemplate<String, Object> redisTemplate;

  @PostConstruct
  public void init() {
    redisTemplate.setValueSerializer(RedisSerializer.java());
    RedisUtils.initialize(redisTemplate);
  }

  @Bean
  public MessageSource iamMessageSource() {
    ReloadableResourceBundleMessageSource messageSource =
        new ReloadableResourceBundleMessageSource();
    messageSource.addBasenames("classpath:i18n/errors/messages");
    messageSource.addBasenames("classpath:i18n/errors/messages-iam");
    messageSource.addBasenames("classpath:i18n/errors/messages-common");
    messageSource.setDefaultLocale(Locale.CHINA);
    return messageSource;
  }

  @Override
  public void addArgumentResolvers(List<HandlerMethodArgumentResolver> resolvers) {
    resolvers.add(new UnderlineToCamelArgumentResolver());
  }
}
