package top.wecoding.iam.sdk;

import cn.hutool.core.util.ReUtil;
import cn.hutool.extra.spring.SpringUtil;
import lombok.Getter;
import lombok.Setter;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.InitializingBean;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.core.annotation.AnnotationUtils;
import org.springframework.web.method.HandlerMethod;
import org.springframework.web.servlet.mvc.method.RequestMappingInfo;
import org.springframework.web.servlet.mvc.method.annotation.RequestMappingHandlerMapping;

import java.util.*;
import java.util.regex.Pattern;

/**
 * 放行白名单配置
 *
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@Slf4j
@ConfigurationProperties(prefix = "security.oauth2.ignore")
public class IgnoreWhiteProperties implements InitializingBean {

  private static final Pattern PATTERN = Pattern.compile("\\{(.*?)}");

  private static final String[] DEFAULT_IGNORE_URLS = {"/actuator/**", "/error", "/v3/api-docs"};

  @Getter @Setter private List<String> whites = new ArrayList<>();

  @Override
  public void afterPropertiesSet() {
    whites.addAll(Arrays.asList(DEFAULT_IGNORE_URLS));

    RequestMappingHandlerMapping mapping = SpringUtil.getBean(RequestMappingHandlerMapping.class);
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
                                .forEach(url -> whites.add(ReUtil.replaceAll(url, PATTERN, "*"))));

                InnerAuth controller =
                    AnnotationUtils.findAnnotation(handlerMethod.getBeanType(), InnerAuth.class);
                Optional.ofNullable(controller)
                    .ifPresent(
                        inner ->
                            Objects.requireNonNull(info.getPatternValues())
                                .forEach(url -> whites.add(ReUtil.replaceAll(url, PATTERN, "*"))));
              } catch (Exception e) {
                log.error("Failed to configure the whitelist: ", e);
              }
            });
  }
}
