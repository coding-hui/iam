package top.wecoding.iam.sdk.interceptor;

import cn.hutool.json.JSONUtil;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.web.servlet.AsyncHandlerInterceptor;
import top.wecoding.iam.sdk.context.AuthContextHolder;
import top.wecoding.iam.common.model.AuthInfo;
import top.wecoding.core.enums.rest.SystemErrorCodeEnum;
import top.wecoding.core.util.ServletResponseUtil;
import top.wecoding.core.util.WebUtil;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

/**
 * 认证拦截器校验
 *
 * @author liuyuhui
 * @date 2022/9/11
 * @qq 1515418211
 */
@Slf4j
public class SecurityInterceptor implements AsyncHandlerInterceptor {

  @Override
  public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler)
      throws Exception {
    AuthInfo authInfo = AuthContextHolder.getContext();
    if (authInfo != null) {
      return true;
    }
    log.warn(
        "Signature authentication fails, request url: {}, ip: {}, params: {}",
        request.getRequestURI(),
        WebUtil.getIP(request),
        JSONUtil.toJsonStr(request.getParameterMap()));
    ServletResponseUtil.webMvcResponseWriter(
        response, HttpStatus.UNAUTHORIZED, SystemErrorCodeEnum.UNAUTHORIZED);
    return false;
  }
}
