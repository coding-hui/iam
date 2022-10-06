package top.wecoding.iam.server.controller;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import top.wecoding.core.result.PageInfo;
import top.wecoding.core.result.R;
import top.wecoding.iam.common.model.request.TokenInfoPageRequest;
import top.wecoding.iam.common.model.response.TokenInfoResponse;
import top.wecoding.iam.server.service.TokenService;

/**
 * @author liuyuhui
 * @date 2022/10/4
 * @qq 1515418211
 */
@Slf4j
@RestController
@RequiredArgsConstructor
@RequestMapping("/api/v1/token")
public class TokenController {

  private final TokenService tokenService;

  @GetMapping("")
  public R<PageInfo<TokenInfoResponse>> page(TokenInfoPageRequest tokenInfoPageRequest) {
    return R.ok(tokenService.infoPage(tokenInfoPageRequest));
  }
}
