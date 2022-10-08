package top.wecoding.iam.sdk;

import java.lang.annotation.*;
import org.springframework.context.annotation.Import;
import org.springframework.security.config.annotation.method.configuration.EnableGlobalMethodSecurity;
import top.wecoding.iam.sdk.feign.IAMFeignRequestInterceptor;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Documented
@Inherited
@Target({ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
@EnableGlobalMethodSecurity(prePostEnabled = true)
@Import({
  IAMResourceServerAutoConfiguration.class,
  IAMResourceServerConfiguration.class,
  IAMFeignRequestInterceptor.class
})
public @interface EnableIAMResourceServer {}
