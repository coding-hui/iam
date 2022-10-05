package top.wecoding.iam.server.util;

import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.junit.jupiter.MockitoExtension;
import org.mockito.junit.jupiter.MockitoSettings;
import org.mockito.quality.Strictness;

/**
 * @author liuyuhui
 * @date 2022/10/1
 * @qq 1515418211
 */
@ExtendWith(MockitoExtension.class)
@MockitoSettings(strictness = Strictness.LENIENT)
public class PasswordUtilsTest {

  @Test
  void checkContent() {
    Assertions.assertTrue(PasswordUtil.checkContent("WeCoding@2022"));
    Assertions.assertFalse(PasswordUtil.checkContent("WeCoding2022"));
    Assertions.assertFalse(PasswordUtil.checkContent("WeCoding"));
    Assertions.assertFalse(PasswordUtil.checkContent("wecoding"));
    Assertions.assertFalse(PasswordUtil.checkContent("2022"));
  }
}
