package top.wecoding.iam.server.config;

import lombok.RequiredArgsConstructor;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.data.redis.serializer.RedisSerializer;
import top.wecoding.iam.server.props.AppProperties;
import top.wecoding.redis.util.RedisUtils;

import javax.annotation.PostConstruct;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Configuration
@RequiredArgsConstructor
@EnableConfigurationProperties(AppProperties.class)
public class IAMServiceConfiguration {

  private final RedisTemplate<String, Object> redisTemplate;

  @PostConstruct
  public void init() {
    redisTemplate.setValueSerializer(RedisSerializer.java());
    RedisUtils.initialize(redisTemplate);
  }
}
