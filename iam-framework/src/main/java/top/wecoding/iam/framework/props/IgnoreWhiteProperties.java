package top.wecoding.iam.framework.props;

import java.util.*;
import java.util.regex.Pattern;
import lombok.Getter;
import lombok.Setter;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.BeansException;
import org.springframework.beans.factory.InitializingBean;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ApplicationContextAware;
import org.springframework.core.annotation.AnnotationUtils;
import org.springframework.stereotype.Component;
import org.springframework.web.method.HandlerMethod;
import org.springframework.web.servlet.mvc.method.RequestMappingInfo;
import org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerMapping;
import top.wecoding.iam.framework.InnerAuth;

/**
 * 放行白名单配置
 *
 * @author liuyuhui
 * @date 2022/10/3
 */
@Slf4j
@Component
public class IgnoreWhiteProperties implements InitializingBean, ApplicationContextAware {

  private static final Pattern PATTERN = Pattern.compile("\\{(.*?)}");

  private static final String[] DEFAULT_IGNORE_URLS = {
    "/actuator/**", "/error", "/auth/*", "/v3/api-docs", "/login"
  };

  private static ApplicationContext context;

  @Getter @Setter private List<String> whites = new ArrayList<>();

  @Override
  public void afterPropertiesSet() {
    whites.addAll(Arrays.asList(DEFAULT_IGNORE_URLS));

    RequestMappingHandlerMapping mapping =
        (RequestMappingHandlerMapping) context.getBean("requestMappingHandlerMapping");
    Map<RequestMappingInfo, HandlerMethod> map = mapping.getHandlerMethods();

    map.keySet()
        .forEach(
            info -> {
              HandlerMethod handlerMethod = map.get(info);
              try {
                InnerAuth method =
                    AnnotationUtils.findAnnotation(handlerMethod.getMethod(), InnerAuth.class);
                Optional.ofNullable(method)
                    .ifPresent(
                        inner ->
                            Objects.requireNonNull(info.getPatternValues())
                                .forEach(
                                    url -> whites.add(url.replaceAll(PATTERN.pattern(), "*"))));

                InnerAuth controller =
                    AnnotationUtils.findAnnotation(handlerMethod.getBeanType(), InnerAuth.class);
                Optional.ofNullable(controller)
                    .ifPresent(
                        inner ->
                            Objects.requireNonNull(info.getPatternValues())
                                .forEach(
                                    url -> whites.add(url.replaceAll(PATTERN.pattern(), "*"))));
              } catch (Exception e) {
                log.error("Failed to configure the whitelist: ", e);
              }
            });
  }

  @Override
  public void setApplicationContext(ApplicationContext applicationContext) throws BeansException {
    IgnoreWhiteProperties.context = applicationContext;
  }
}
