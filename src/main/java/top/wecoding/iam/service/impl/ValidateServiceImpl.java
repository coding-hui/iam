package top.wecoding.iam.service.impl;

import org.springframework.stereotype.Service;
import top.wecoding.core.result.R;
import top.wecoding.iam.service.ValidateService;

/**
 * @author liuyuhui
 * @qq 1515418211
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
