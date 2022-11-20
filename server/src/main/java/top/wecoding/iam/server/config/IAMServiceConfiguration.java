package top.wecoding.iam.server.config;

import lombok.RequiredArgsConstructor;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.data.redis.serializer.RedisSerializer;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.web.method.support.HandlerMethodArgumentResolver;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
import top.wecoding.iam.server.props.AppProperties;
import top.wecoding.iam.server.web.UnderlineToCamelArgumentResolver;
import top.wecoding.redis.util.RedisUtils;

import javax.annotation.PostConstruct;
import java.util.List;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@EnableWebSecurity
@RequiredArgsConstructor
@Configuration(proxyBeanMethods = false)
@EnableConfigurationProperties(AppProperties.class)
public class IAMServiceConfiguration implements WebMvcConfigurer {

  private final RedisTemplate<String, Object> redisTemplate;

  @PostConstruct
  public void init() {
    redisTemplate.setValueSerializer(RedisSerializer.java());
    RedisUtils.initialize(redisTemplate);
  }

  @Override
  public void addArgumentResolvers(List<HandlerMethodArgumentResolver> resolvers) {
    resolvers.add(new UnderlineToCamelArgumentResolver());
  }
}
