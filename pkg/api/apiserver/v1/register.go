// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import "github.com/coding-hui/common/scheme"

// GroupName is the group name use in this package.
// If use a public domain name, need set the GroupName to service name.
// For example: if restful path is: https://wecoding/apimachinery/v1/secrets, we can set GroupName="apimachinery".
const GroupName = "api"

// SchemeGroupVersion is group version used to register these objects.
var SchemeGroupVersion = scheme.GroupVersion{Group: GroupName, Version: "v1"}

// Resource takes an unqualified resource and returns a Group qualified GroupResource.
func Resource(resource string) scheme.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}
