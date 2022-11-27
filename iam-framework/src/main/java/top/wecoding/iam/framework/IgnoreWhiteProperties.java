package top.wecoding.iam.framework;

import lombok.Getter;
import lombok.Setter;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.InitializingBean;
import org.springframework.stereotype.Component;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.regex.Pattern;

/**
 * 放行白名单配置
 *
 * @author liuyuhui
 * @date 2022/10/3
 */
@Slf4j
@Component
public class IgnoreWhiteProperties implements InitializingBean {

  private static final Pattern PATTERN = Pattern.compile("\\{(.*?)}");

  private static final String[] DEFAULT_IGNORE_URLS = {"/actuator/**", "/error", "/v3/api-docs"};

  @Getter @Setter private List<String> whites = new ArrayList<>();

  @Override
  public void afterPropertiesSet() {
    whites.addAll(Arrays.asList(DEFAULT_IGNORE_URLS));
    //
    // RequestMappingHandlerMapping mapping = SpringUtil.getBean("requestMappingHandlerMapping");
    // Map<RequestMappingInfo, HandlerMethod> map = mapping.getHandlerMethods();
    //
    // map.keySet()
    //     .forEach(
    //         info -> {
    //           HandlerMethod handlerMethod = map.get(info);
    //           try {
    //             InnerAuth method =
    //                 AnnotationUtils.findAnnotation(handlerMethod.getMethod(), InnerAuth.class);
    //             Optional.ofNullable(method)
    //                 .ifPresent(
    //                     inner ->
    //                         Objects.requireNonNull(info.getPatternValues())
    //                             .forEach(url -> whites.add(ReUtil.replaceAll(url, PATTERN,
    // "*"))));
    //
    //             InnerAuth controller =
    //                 AnnotationUtils.findAnnotation(handlerMethod.getBeanType(), InnerAuth.class);
    //             Optional.ofNullable(controller)
    //                 .ifPresent(
    //                     inner ->
    //                         Objects.requireNonNull(info.getPatternValues())
    //                             .forEach(url -> whites.add(ReUtil.replaceAll(url, PATTERN,
    // "*"))));
    //           } catch (Exception e) {
    //             log.error("Failed to configure the whitelist: ", e);
    //           }
    //         });
  }
}
