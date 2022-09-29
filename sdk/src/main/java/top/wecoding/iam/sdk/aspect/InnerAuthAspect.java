package top.wecoding.iam.sdk.aspect;

import cn.hutool.core.util.StrUtil;
import lombok.extern.slf4j.Slf4j;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.springframework.core.Ordered;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.core.exception.user.InnerAuthException;
import top.wecoding.core.util.WebUtil;
import top.wecoding.iam.sdk.annotation.InnerAuth;

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
    String source = WebUtil.getRequest().getHeader(SecurityConstants.FROM_SOURCE);

    // 是否是内部请求
    if (!StrUtil.equals(SecurityConstants.INNER, source)) {
      log.warn("访问接口 {} 没有权限", point.getSignature().getName());
      throw new InnerAuthException();
    }

    // 用户验证
    String userId = WebUtil.getRequest().getHeader(SecurityConstants.DETAILS_USER_ID);
    String account = WebUtil.getRequest().getHeader(SecurityConstants.DETAILS_ACCOUNT);
    if (innerAuth.value() && StrUtil.hasBlank(userId, account)) {
      log.warn("访问接口 {} 没有权限", point.getSignature().getName());
      throw new InnerAuthException();
    }

    log.debug(" >>> 内部应用访问，不检验权限，访问接口: {}", point.getSignature().getName());
    return point.proceed();
  }

  @Override
  public int getOrder() {
    return Ordered.HIGHEST_PRECEDENCE + 1;
  }
}
