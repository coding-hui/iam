package top.wecoding.iam.server.web.annotation;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * @author liuyuhui
 * @date 2022/11/5
 */
@Target(value = ElementType.PARAMETER)
@Retention(RetentionPolicy.RUNTIME)
public @interface RequestParameter {}
