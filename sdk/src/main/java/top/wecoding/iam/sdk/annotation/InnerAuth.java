package top.wecoding.iam.sdk.annotation;

import org.springframework.core.annotation.AliasFor;

import java.lang.annotation.*;

/**
 * 内部服务认证
 *
 * @author liuyuhui
 * @qq 1515418211
 */
@Target(ElementType.METHOD)
@Retention(RetentionPolicy.RUNTIME)
@Documented
public @interface InnerAuth {

  /** 是否进行用户信息校验 */
  @AliasFor("requiresLogin")
  boolean value() default false;

  /** 是否进行用户信息校验 */
  @AliasFor("value")
  boolean requiresLogin() default false;
}
