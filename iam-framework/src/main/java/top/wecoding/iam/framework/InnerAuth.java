package top.wecoding.iam.framework;

import java.lang.annotation.Documented;
import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * 内部服务认证
 *
 * @author liuyuhui
 */
@Target({ElementType.METHOD, ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
@Documented
public @interface InnerAuth {

  /** 是否校验请求头来源，只允许内部访问 */
  boolean value() default true;
}
