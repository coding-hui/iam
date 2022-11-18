package top.wecoding.iam.framework;

import java.lang.annotation.*;

/**
 * 内部服务认证
 *
 * @author liuyuhui
 * @qq 1515418211
 */
@Target({ElementType.METHOD, ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
@Documented
public @interface InnerAuth {

  /** 是否校验请求头来源，只允许内部访问 */
  boolean value() default true;
}
