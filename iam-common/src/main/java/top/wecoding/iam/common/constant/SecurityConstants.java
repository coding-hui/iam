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
 */
public interface SecurityConstants {

  /** 项目的license */
  String PROJECT_LICENSE = "http://localhost";

  /** 协议字段 */
  String DETAILS_LICENSE = "license";

  /** 客户端ID */
  String CLIENT_ID = "client_id";

  /** 客户端模式 */
  String CLIENT_CREDENTIALS = "client_credentials";

  /** 用户信息 */
  String DETAILS_USER = "user_info";

  /** 请求来源 */
  String FROM = "from-source";

  /** 内部请求 */
  String INNER = "inner";

  /** 手机号登录 */
  String APP = "app";

  /** {bcrypt} 加密的特征码 */
  String BCRYPT = "{bcrypt}";

  /** {noop} 加密的特征码 */
  String NOOP = "{noop}";

  /** 授权码模式confirm */
  String CUSTOM_CONSENT_PAGE_URI = "/token/confirm_access";

  /** The default endpoint {@code URI} for access token requests. */
  String CUSTOM_TOKEN_ENDPOINT_URI = "/oauth2/token";

  /** role prefix */
  String ROLE_PREFIX = "ROLE_";
}
