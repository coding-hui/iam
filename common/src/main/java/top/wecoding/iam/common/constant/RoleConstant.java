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
package top.wecoding.iam.common.constant;

/**
 * @author liuyuhui
 * @date 2022
 * @qq 1515418211
 */
public class RoleConstant {

  public static final String ALL_PERMISSION = "*:*:*";

  public static final String ADMIN = "admin";

  public static final String HAS_ROLE_ADMIN = "hasAnyRole('" + ADMIN + "')";

  public static final String USER = "user";

  public static final String HAS_ROLE_USER = "hasAnyRole('" + USER + "')";

  public static final String TEST = "test";

  public static final String HAS_ROLE_TEST = "hasAnyRole('" + TEST + "')";
}
