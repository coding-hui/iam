package top.wecoding.iam.server.util;

import java.util.Arrays;
import java.util.Map;
import java.util.function.Function;
import java.util.function.Supplier;
import java.util.stream.Collectors;
import lombok.experimental.UtilityClass;

/**
 * @author liuyuhui
 * @qq 1515418211
 */
@UtilityClass
public class EnumUtil {

  public <T> Map<String, T> ofDict(Supplier<T[]> supplier, Function<T, Object> function) {
    return Arrays.stream(supplier.get())
        .collect(Collectors.toMap(t -> function.apply(t).toString(), Function.identity()));
  }
}
