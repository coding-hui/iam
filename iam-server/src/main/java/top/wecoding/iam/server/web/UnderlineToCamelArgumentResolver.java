package top.wecoding.iam.server.web;

import java.util.Iterator;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import lombok.extern.slf4j.Slf4j;
import org.jetbrains.annotations.NotNull;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.BeanWrapper;
import org.springframework.beans.PropertyAccessorFactory;
import org.springframework.core.MethodParameter;
import org.springframework.web.bind.support.WebDataBinderFactory;
import org.springframework.web.context.request.NativeWebRequest;
import org.springframework.web.method.support.ModelAndViewContainer;
import top.wecoding.iam.server.web.annotation.RequestParameter;

/**
 * @author liuyuhui
 * @date 2022/11/5
 */
@Slf4j
public class UnderlineToCamelArgumentResolver extends AbstractCustomizeResolver {

  private static final Pattern UNDER_LINE_PATTERN = Pattern.compile("_(\\w)");

  @Override
  public boolean supportsParameter(MethodParameter parameter) {
    return parameter.hasParameterAnnotation(RequestParameter.class);
  }

  @Override
  public Object resolveArgument(
      @NotNull MethodParameter parameter,
      ModelAndViewContainer mavContainer,
      @NotNull NativeWebRequest webRequest,
      WebDataBinderFactory binderFactory)
      throws Exception {
    Object org = handleParameterNames(parameter, webRequest);
    valid(parameter, mavContainer, webRequest, binderFactory, org);
    return org;
  }

  private Object handleParameterNames(MethodParameter parameter, NativeWebRequest webRequest) {
    Object obj = BeanUtils.instantiateClass(parameter.getParameterType());
    BeanWrapper wrapper = PropertyAccessorFactory.forBeanPropertyAccess(obj);
    Iterator<String> paramNames = webRequest.getParameterNames();
    while (paramNames.hasNext()) {
      String paramName = paramNames.next();
      Object o = webRequest.getParameter(paramName);
      try {
        wrapper.setPropertyValue(underLineToCamel(paramName), o);
      } catch (Exception e) {
        log.warn("Failed to set the attribute value: {}, error: {}", parameter, e.getMessage());
      }
    }
    return obj;
  }

  private String underLineToCamel(String source) {
    Matcher matcher = UNDER_LINE_PATTERN.matcher(source);
    StringBuffer result = new StringBuffer();
    while (matcher.find()) {
      matcher.appendReplacement(result, matcher.group(1).toUpperCase());
    }
    matcher.appendTail(result);
    return result.toString();
  }
}
