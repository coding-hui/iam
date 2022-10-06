package top.wecoding.iam.server.config;

import org.springframework.boot.autoconfigure.cache.CacheProperties;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.data.redis.serializer.StringRedisSerializer;
import top.wecoding.iam.server.props.AppProperties;
import top.wecoding.redis.repository.RedisCacheOperatorImpl;
import top.wecoding.redis.repository.base.CacheOperatorPlus;
import top.wecoding.redis.serializer.FastJson2JsonRedisSerializer;
import top.wecoding.redis.service.RedisService;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Configuration
@EnableConfigurationProperties(AppProperties.class)
public class IAMServiceConfiguration {

  @Bean
  public RedisService redisService(
      RedisTemplate<String, Object> redisTemplate,
      StringRedisTemplate stringRedisTemplate,
      CacheProperties properties) {

    FastJson2JsonRedisSerializer serializer = new FastJson2JsonRedisSerializer(Object.class);

    // 使用StringRedisSerializer来序列化和反序列化redis的key值
    redisTemplate.setKeySerializer(new StringRedisSerializer());
    redisTemplate.setValueSerializer(serializer);

    // Hash的key也采用StringRedisSerializer的序列化方式
    redisTemplate.setHashKeySerializer(new StringRedisSerializer());
    redisTemplate.setHashValueSerializer(serializer);
    return new RedisService(
        redisTemplate, stringRedisTemplate, properties.getRedis().isCacheNullValues());
  }

  @Bean
  public CacheOperatorPlus cacheOperator(RedisService redisService) {
    return new RedisCacheOperatorImpl(redisService);
  }
}
