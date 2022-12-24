package top.wecoding.iam.framework.security.web;

import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import java.io.PrintWriter;
import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.oauth2.server.resource.InvalidBearerTokenException;
import org.springframework.security.web.AuthenticationEntryPoint;
import org.springframework.stereotype.Component;
import top.wecoding.commons.core.constant.StrPool;
import top.wecoding.commons.core.enums.SystemErrorCodeEnum;
import top.wecoding.commons.core.model.R;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
@Component
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
    response.setStatus(HttpStatus.UNAUTHORIZED.value());
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
