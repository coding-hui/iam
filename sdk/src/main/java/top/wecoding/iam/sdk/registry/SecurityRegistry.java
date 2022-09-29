package top.wecoding.iam.sdk.registry;

import lombok.Data;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

/**
 * API 放行配置
 *
 * @author liuyuhui
 * @date 2022/6/6
 * @qq 1515418211
 */
@Data
public class SecurityRegistry {

  private List<String> defaultExcludePatterns = new ArrayList<>();
  private List<String> excludePatterns = new ArrayList<>();
  private List<String> includePatterns = new ArrayList<>();
  private boolean enabled = true;
  private int order = 0;

  public SecurityRegistry() {
    this.defaultExcludePatterns.add("/token/**");
    this.defaultExcludePatterns.add("/auth/**");
    this.defaultExcludePatterns.add("/*/v2/api-docs");
    this.defaultExcludePatterns.add("/actuator/health/**");
    this.includePatterns.add("/*");
  }

  public SecurityRegistry excludePathPatterns(String... patterns) {
    return excludePathPatterns(Arrays.asList(patterns));
  }

  public SecurityRegistry excludePathPatterns(List<String> patterns) {
    this.excludePatterns.addAll(patterns);
    return this;
  }

  public SecurityRegistry addPathPatterns(String... patterns) {
    return addPathPatterns(Arrays.asList(patterns));
  }

  public SecurityRegistry addPathPatterns(List<String> patterns) {
    this.includePatterns =
        (this.includePatterns != null ? this.includePatterns : new ArrayList<>(patterns.size()));
    this.includePatterns.addAll(patterns);
    return this;
  }
}
