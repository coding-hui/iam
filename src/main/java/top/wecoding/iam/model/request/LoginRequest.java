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
package top.wecoding.iam.model.request;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import lombok.*;
import lombok.experimental.Accessors;

import javax.validation.constraints.NotEmpty;

/**
 * 登录参数
 *
 * @author liuyuhui
 * @date 2022
 * @qq 1515418211
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Accessors(chain = true)
@EqualsAndHashCode(callSuper = false)
@ApiModel(value = "LoginParamDTO", description = "登录参数")
public class LoginRequest {

  @ApiModelProperty(value = "账号")
  private String account;

  @ApiModelProperty(value = "密码")
  private String password;

  @ApiModelProperty(value = "刷新令牌")
  private String refreshToken;

  /** password: 账号密码 captcha: 验证码 */
  @ApiModelProperty(value = "授权类型", example = "password", allowableValues = "captcha,password")
  @NotEmpty(message = "授权类型不能为空")
  private String grantType;
}
