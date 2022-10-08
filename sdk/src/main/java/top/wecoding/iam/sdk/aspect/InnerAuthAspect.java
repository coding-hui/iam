package top.wecoding.iam.sdk.aspect;

import cn.hutool.core.util.StrUtil;
import javax.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.AnnotationUtils;
import org.springframework.security.access.AccessDeniedException;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.core.util.WebUtils;
import top.wecoding.iam.sdk.InnerAuth;

/**
 * 内部服务调用验证处理切面
 *
 * @author liuyuhui
 * @qq 1515418211
 */
@Slf4j
@Aspect
public class InnerAuthAspect implements Ordered {

  @Around("@within(innerAuth) || @annotation(innerAuth)")
  public Object innerAround(ProceedingJoinPoint point, InnerAuth innerAuth) throws Throwable {
    if (innerAuth == null) {
      Class<?> clazz = point.getTarget().getClass();
      innerAuth = AnnotationUtils.findAnnotation(clazz, InnerAuth.class);
    }

    HttpServletRequest request = WebUtils.getRequest();
    String requestURI = request.getRequestURI();
    String source = request.getHeader(SecurityConstants.FROM);
    if (innerAuth != null
        && innerAuth.value()
        && !StrUtil.equals(SecurityConstants.INNER, source)) {
      log.warn("Access api {} does not have permission", requestURI);
      throw new AccessDeniedException("Access is denied");
    }

    log.debug("Internal request, skip authentication: {}", requestURI);
    return point.proceed();
  }

  @Override
  public int getOrder() {
    return Ordered.HIGHEST_PRECEDENCE + 1;
  }
}
