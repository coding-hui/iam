package top.wecoding.iam.framework.web;

import cn.hutool.http.HttpStatus;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import org.springframework.http.MediaType;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.oauth2.server.resource.InvalidBearerTokenException;
import org.springframework.security.web.AuthenticationEntryPoint;
import top.wecoding.core.constant.StrPool;
import top.wecoding.core.enums.rest.SystemErrorCodeEnum;
import top.wecoding.core.result.R;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.PrintWriter;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
@RequiredArgsConstructor
public class ResourceAuthExceptionEntryPoint implements AuthenticationEntryPoint {

  private final ObjectMapper objectMapper;

  @Override
  @SneakyThrows
  public void commence(
      HttpServletRequest request,
      HttpServletResponse response,
      AuthenticationException authException) {
    response.setCharacterEncoding(StrPool.UTF8);
    response.setContentType(MediaType.APPLICATION_JSON_VALUE);
    R<String> result = new R<>();
    result.setCode(SystemErrorCodeEnum.FAILURE.getCode());
    response.setStatus(HttpStatus.HTTP_UNAUTHORIZED);
    if (authException != null) {
      result.setMsg("error");
      result.setData(authException.getMessage());
    }

    // 针对令牌过期返回特殊的 424
    if (authException instanceof InvalidBearerTokenException) {
      response.setStatus(org.springframework.http.HttpStatus.FAILED_DEPENDENCY.value());
      result.setMsg("token expire");
    }
    PrintWriter printWriter = response.getWriter();
    printWriter.append(objectMapper.writeValueAsString(result));
  }
}
