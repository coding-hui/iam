package top.wecoding.iam.server.config;

import java.util.List;
import javax.annotation.PostConstruct;
import lombok.RequiredArgsConstructor;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.data.redis.serializer.RedisSerializer;
import org.springframework.web.method.support.HandlerMethodArgumentResolver;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;
import top.wecoding.iam.server.props.AppProperties;
import top.wecoding.iam.server.web.UnderlineToCamelArgumentResolver;
import top.wecoding.redis.util.RedisUtils;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Configuration
@RequiredArgsConstructor
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
