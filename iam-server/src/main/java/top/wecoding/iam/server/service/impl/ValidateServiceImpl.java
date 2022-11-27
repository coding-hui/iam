package top.wecoding.iam.server.service.impl;

import org.springframework.stereotype.Service;
import top.wecoding.commons.core.model.R;
import top.wecoding.iam.server.service.ValidateService;

/**
 * @author liuyuhui
 */
@Service
public class ValidateServiceImpl implements ValidateService {

  @Override
  public R<?> createCode() {
    return R.ok("1234");
  }

  @Override
  public R<?> createSmsCode(String mobile) {
    return R.ok(mobile.substring(0, 4));
  }

  @Override
  public void check(String key, String code) {}
}
