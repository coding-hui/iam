package top.wecoding.iam.server.advice;

import lombok.extern.slf4j.Slf4j;
import org.springframework.context.MessageSource;
import org.springframework.context.i18n.LocaleContextHolder;
import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestControllerAdvice;
import top.wecoding.commons.core.exception.handler.AbstractGlobalExceptionHandler;
import top.wecoding.commons.core.exception.user.UnauthorizedException;
import top.wecoding.commons.core.model.R;
import top.wecoding.iam.common.enums.IamErrorCode;

/**
 * @author liuyuhui
 * @date 2022/10/2
 */
@Slf4j
@RestControllerAdvice
@Order(Ordered.HIGHEST_PRECEDENCE)
public class WeCodingExceptionAdvice extends AbstractGlobalExceptionHandler {

  public WeCodingExceptionAdvice(MessageSource iamMessageSource) {
    super(iamMessageSource);
  }

  @ResponseStatus(HttpStatus.UNAUTHORIZED)
  @ExceptionHandler(UnauthorizedException.class)
  public R<Object> unauthorizedExceptionHandle(UnauthorizedException e) {
    log.warn("UnauthorizedException: {}", e.getMessage());
    String message =
        messageSource.getMessage(
            IamErrorCode.UNAUTHORIZED.getCode(),
            new Object[] {e.getMessage()},
            LocaleContextHolder.getLocale());
    return R.error(IamErrorCode.UNAUTHORIZED, message);
  }
}
