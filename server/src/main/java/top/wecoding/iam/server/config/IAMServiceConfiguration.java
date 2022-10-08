package top.wecoding.iam.server.config;

import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Configuration;
import top.wecoding.iam.server.props.AppProperties;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Configuration
@EnableConfigurationProperties(AppProperties.class)
public class IAMServiceConfiguration {
}
