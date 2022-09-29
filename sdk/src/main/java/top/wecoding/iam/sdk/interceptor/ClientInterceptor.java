package top.wecoding.iam.sdk.interceptor;

import cn.hutool.core.util.StrUtil;
import cn.hutool.json.JSONUtil;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.web.servlet.AsyncHandlerInterceptor;
import top.wecoding.core.enums.rest.SystemErrorCodeEnum;
import top.wecoding.core.util.ServletResponseUtil;
import top.wecoding.core.util.WebUtil;
import top.wecoding.iam.common.helper.SecurityHelper;
import top.wecoding.iam.common.model.AuthInfo;
import top.wecoding.iam.sdk.context.AuthContextHolder;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

/**
 * 客户端校验
 *
 * @author liuyuhui
 * @date 2022/9/11
 * @qq 1515418211
 */
@Slf4j
@AllArgsConstructor
public class ClientInterceptor implements AsyncHandlerInterceptor {

  private final String clientId;

  @Override
  public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler)
      throws Exception {
    AuthInfo authInfo = AuthContextHolder.getContext();
    if (authInfo != null
        && StrUtil.equals(clientId, SecurityHelper.getClientIdFromHeader())
        && StrUtil.equals(clientId, authInfo.getClientId())) {
      return true;
    } else {
      log.warn(
          "Client authentication fails, request url: {}, ip: {}, params: {}",
          request.getRequestURI(),
          WebUtil.getIP(request),
          JSONUtil.toJsonStr(request.getParameterMap()));
      ServletResponseUtil.webMvcResponseWriter(
          response, HttpStatus.OK, SystemErrorCodeEnum.UNAUTHORIZED);
      return false;
    }
  }
}
