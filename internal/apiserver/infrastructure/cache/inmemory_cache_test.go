// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cache

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/coding-hui/common/util/sets"
)

var dataSet = map[string]string{
	"foo1": "val1",
	"foo2": "val2",
	"foo3": "val3",
	"bar1": "val1",
	"bar2": "val2",
}

// load dataset into cache
func load(client Interface, data map[string]string) error {
	for k, v := range data {
		err := client.Set(k, v, NeverExpire)
		if err != nil {
			return err
		}
	}

	return nil
}

// dump retrieve all data in simple into a map
func dump(client Interface) (map[string]string, error) {
	keys, err := client.Keys("*")
	if err != nil {
		return nil, err
	}

	snapshot := make(map[string]string)
	for _, key := range keys {
		val, err := client.Get(key)
		if err != nil {
			continue
		}
		snapshot[key] = val
	}

	return snapshot, nil
}

func TestDeleteAndExpireCache(t *testing.T) {
	var testCases = []struct {
		description    string
		deleteKeys     sets.String
		expireKeys     sets.String
		expireDuration time.Duration // never use a 0(NeverExpires) duration with expireKeys, recommend time.Millisecond * 500.
		expected       map[string]string
	}{
		{
			description: "Should get all keys",
			expected: map[string]string{
				"foo1": "val1",
				"foo2": "val2",
				"foo3": "val3",
				"bar1": "val1",
				"bar2": "val2",
			},
		},
		{
			description: "Test delete should get only keys start with foo",
			expected: map[string]string{
				"foo1": "val1",
				"foo2": "val2",
				"foo3": "val3",
			},
			deleteKeys: sets.NewString("bar1", "bar2"),
		},
		{
			description: "Should get only keys start with bar",
			expected: map[string]string{
				"bar1": "val1",
				"bar2": "val2",
			},
			expireDuration: time.Millisecond * 500,
			expireKeys:     sets.NewString("foo1", "foo2", "foo3"),
		},
	}

	for _, testCase := range testCases {
		cacheClient, _ := NewInMemoryCache(nil, nil)

		t.Run(testCase.description, func(t *testing.T) {
			err := load(cacheClient, dataSet)
			if err != nil {
				t.Fatalf("Unable to load dataset, got error %v", err)
			}

			if len(testCase.deleteKeys) != 0 {
				err = cacheClient.Del(testCase.deleteKeys.UnsortedList()...)
				if err != nil {
					t.Fatalf("Error delete keys, %v", err)
				}
			}

			if len(testCase.expireKeys) != 0 && testCase.expireDuration != 0 {
				for _, key := range testCase.expireKeys.UnsortedList() {
					err = cacheClient.Expire(key, testCase.expireDuration)
					if err != nil {
						t.Fatalf("Error expire keys, %v", err)
					}
				}
				time.Sleep(testCase.expireDuration)
			}

			got, err := dump(cacheClient)
			if err != nil {
				t.Fatalf("Error dump data, %v", err)
			}

			if diff := cmp.Diff(got, testCase.expected); len(diff) != 0 {
				t.Errorf("%T differ (-got, +expected) %v", testCase.expected, diff)
			}
		})
	}
}
