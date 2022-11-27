package top.wecoding.iam.server.web;

import org.springframework.core.Conventions;
import org.springframework.core.MethodParameter;
import org.springframework.core.annotation.AnnotationUtils;
import org.springframework.validation.BindingResult;
import org.springframework.validation.Errors;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.WebDataBinder;
import org.springframework.web.bind.support.WebDataBinderFactory;
import org.springframework.web.context.request.NativeWebRequest;
import org.springframework.web.method.support.HandlerMethodArgumentResolver;
import org.springframework.web.method.support.ModelAndViewContainer;

import java.lang.annotation.Annotation;
import java.util.Objects;

/**
 * @author liuyuhui
 * @date 2022/11/5
 */
public abstract class AbstractCustomizeResolver implements HandlerMethodArgumentResolver {

  protected void valid(
      MethodParameter parameter,
      ModelAndViewContainer mavContainer,
      NativeWebRequest webRequest,
      WebDataBinderFactory binderFactory,
      Object arg)
      throws Exception {
    String name = Conventions.getVariableNameForParameter(parameter);
    WebDataBinder binder = binderFactory.createBinder(webRequest, arg, name);
    if (arg != null) {
      validateIfApplicable(binder, parameter);
      if (binder.getBindingResult().hasErrors() && isBindExceptionRequired(binder, parameter)) {
        throw new MethodArgumentNotValidException(parameter, binder.getBindingResult());
      }
    }
    mavContainer.addAttribute(BindingResult.MODEL_KEY_PREFIX + name, binder.getBindingResult());
  }

  protected void validateIfApplicable(WebDataBinder binder, MethodParameter parameter) {
    Annotation[] annotations = parameter.getParameterAnnotations();
    for (Annotation ann : annotations) {
      Validated validatedAnn = AnnotationUtils.getAnnotation(ann, Validated.class);
      if (validatedAnn != null || ann.annotationType().getSimpleName().startsWith("Valid")) {
        Object hints =
            (validatedAnn != null ? validatedAnn.value() : AnnotationUtils.getValue(ann));
        Object[] validationHints =
            (hints instanceof Object[] ? (Object[]) hints : new Object[] {hints});
        binder.validate(validationHints);
        break;
      }
    }
  }

  protected boolean isBindExceptionRequired(WebDataBinder binder, MethodParameter parameter) {
    int i = parameter.getParameterIndex();
    Class<?>[] paramTypes = Objects.requireNonNull(parameter.getMethod()).getParameterTypes();
    boolean hasBindingResult =
        (paramTypes.length > (i + 1) && Errors.class.isAssignableFrom(paramTypes[i + 1]));
    return !hasBindingResult;
  }
}
