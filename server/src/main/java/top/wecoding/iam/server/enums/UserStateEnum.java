/*
 * Copyright (c) 2022. WeCoding (wecoding@yeah.net).
 *
 * Licensed under the GNU LESSER GENERAL PUBLIC LICENSE 3.0;
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.gnu.org/licenses/lgpl.html
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package top.wecoding.iam.server.enums;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.experimental.Accessors;
import top.wecoding.iam.server.util.EnumUtil;

import java.util.Map;

/**
 * 用户类型枚举
 *
 * @author liuyuhui
 * @qq 1515418211
 */
@Getter
@AllArgsConstructor
@Accessors(fluent = true)
public enum UserStateEnum {
  DEFAULT(0),
  INACTIVATED(1),
  DISABLE(2);

  private static final Map<String, UserStateEnum> DICT =
      EnumUtil.ofDict(UserStateEnum::values, UserStateEnum::code);

  private final int code;

  public static UserStateEnum of(int code) {
    return of(String.valueOf(code));
  }

  public static UserStateEnum of(String code) {
    return DICT.get(code);
  }

  public boolean is(int code) {
    return this.code == code;
  }

  public boolean is(String code) {
    return this == DICT.get(code);
  }
}
