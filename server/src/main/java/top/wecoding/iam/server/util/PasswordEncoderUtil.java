package top.wecoding.iam.server.util;

import lombok.experimental.UtilityClass;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.DelegatingPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;

import java.util.Collections;
import java.util.Map;

/**
 * @author liuyuhui
 * @qq 1515418211
 */
@UtilityClass
public class PasswordEncoderUtil {

  private static final String idForEncode = "Bcrypt";

  private static final Map<String, PasswordEncoder> idToPasswordEncoder =
      Collections.singletonMap(idForEncode, new BCryptPasswordEncoder());

  private static DelegatingPasswordEncoder passwordEncoder =
      new DelegatingPasswordEncoder(idForEncode, idToPasswordEncoder);

  public static String encode(String rawPassword) {
    return passwordEncoder.encode(rawPassword);
  }

  public static boolean matches(CharSequence rawPassword, String prefixEncodedPassword) {
    return passwordEncoder.matches(rawPassword, prefixEncodedPassword);
  }
}
