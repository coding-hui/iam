package top.wecoding.iam.sdk.annotation;

import java.lang.annotation.*;

/**
 * 权限检验注解
 *
 * @author liuyuhui
 * @date 2022
 * @qq 1515418211
 */
@Target({ElementType.METHOD, ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
@Inherited
@Documented
public @interface PreAuth {

  /**
   * Spring el 文档地址： <a
   * href="https://docs.spring.io/spring/docs/4.3.16.RELEASE/spring-framework-reference/htmlsingle/#expressions-operators-logical">el</a>
   */
  String value();
}
