/**
 *  Copyright 2015 Paul Querna, 2018 Luis Gomez
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package cacheobject

import (
	"fmt"
	"strings"
	"time"
)

var (
	defaultMaxAge = time.Hour * 24 * 365
)

// TODO(pquerna): add extensions from here: http://www.iana.org/assignments/http-cache-directives/http-cache-directives.xhtml

// DeltaSeconds specifies a non-negative integer, representing
// time in seconds: http://tools.ietf.org/html/rfc7234#section-1.2.1
//
// When set to -1, this means unset.
//
type DeltaSeconds int32

// Fields present in a header.
type FieldNames map[string]bool

// LOW LEVEL API: Repersentation of possible response directives in a `Cache-Control` header: http://tools.ietf.org/html/rfc7234#section-5.2.2
//
// Note: Many fields will be `nil` in practice.
//
type ResponseCacheDirectives struct {

	// must-revalidate(bool): http://tools.ietf.org/html/rfc7234#section-5.2.2.1
	//
	// The "must-revalidate" response directive indicates that once it has
	// become stale, a cache MUST NOT use the response to satisfy subsequent
	// requests without successful validation on the origin server.
	MustRevalidate bool

	// no-cache(FieldName): http://tools.ietf.org/html/rfc7234#section-5.2.2.2
	//
	// The "no-cache" response directive indicates that the response MUST
	// NOT be used to satisfy a subsequent request without successful
	// validation on the origin server.
	//
	// If the no-cache response directive specifies one or more field-names,
	// then a cache MAY use the response to satisfy a subsequent request,
	// subject to any other restrictions on caching.  However, any header
	// fields in the response that have the field-name(s) listed MUST NOT be
	// sent in the response to a subsequent request without successful
	// revalidation with the origin server.
	NoCache FieldNames

	// no-cache(cast-to-bool): http://tools.ietf.org/html/rfc7234#section-5.2.2.2
	//
	// While the RFC defines optional field-names on a no-cache directive,
	// many applications only want to know if any no-cache directives were
	// present at all.
	NoCachePresent bool

	// no-store(bool): http://tools.ietf.org/html/rfc7234#section-5.2.2.3
	//
	// The "no-store" request directive indicates that a cache MUST NOT
	// store any part of either this request or any response to it.  This
	// directive applies to both private and shared caches.
	NoStore bool

	// no-transform(bool): http://tools.ietf.org/html/rfc7234#section-5.2.2.4
	//
	// The "no-transform" response directive indicates that an intermediary
	// (regardless of whether it implements a cache) MUST NOT transform the
	// payload, as defined in Section 5.7.2 of RFC7230.
	NoTransform bool

	// public(bool): http://tools.ietf.org/html/rfc7234#section-5.2.2.5
	//
	// The "public" response directive indicates that any cache MAY store
	// the response, even if the response would normally be non-cacheable or
	// cacheable only within a private cache.
	Public bool

	// private(FieldName): http://tools.ietf.org/html/rfc7234#section-5.2.2.6
	//
	// The "private" response directive indicates that the response message
	// is intended for a single user and MUST NOT be stored by a shared
	// cache.  A private cache MAY store the response and reuse it for later
	// requests, even if the response would normally be non-cacheable.
	//
	// If the private response directive specifies one or more field-names,
	// this requirement is limited to the field-values associated with the
	// listed response header fields.  That is, a shared cache MUST NOT
	// store the specified field-names(s), whereas it MAY store the
	// remainder of the response message.
	Private FieldNames

	// private(cast-to-bool): http://tools.ietf.org/html/rfc7234#section-5.2.2.6
	//
	// While the RFC defines optional field-names on a private directive,
	// many applications only want to know if any private directives were
	// present at all.
	PrivatePresent bool

	// proxy-revalidate(bool): http://tools.ietf.org/html/rfc7234#section-5.2.2.7
	//
	// The "proxy-revalidate" response directive has the same meaning as the
	// must-revalidate response directive, except that it does not apply to
	// private caches.
	ProxyRevalidate bool

	// max-age(delta seconds): http://tools.ietf.org/html/rfc7234#section-5.2.2.8
	//
	// The "max-age" response directive indicates that the response is to be
	// considered stale after its age is greater than the specified number
	// of seconds.
	MaxAge DeltaSeconds

	// s-maxage(delta seconds): http://tools.ietf.org/html/rfc7234#section-5.2.2.9
	//
	// The "s-maxage" response directive indicates that, in shared caches,
	// the maximum age specified by this directive overrides the maximum age
	// specified by either the max-age directive or the Expires header
	// field.  The s-maxage directive also implies the semantics of the
	// proxy-revalidate response directive.
	SMaxAge DeltaSeconds

	////
	// Experimental features
	// - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control#Extension_Cache-Control_directives
	// - https://www.fastly.com/blog/stale-while-revalidate-stale-if-error-available-today
	////

	// immutable(cast-to-bool): experimental feature
	Immutable bool

	// stale-if-error(delta seconds): experimental feature
	StaleIfError DeltaSeconds

	// stale-while-revalidate(delta seconds): experimental feature
	StaleWhileRevalidate DeltaSeconds

	// Extensions: http://tools.ietf.org/html/rfc7234#section-5.2.3
	//
	// The Cache-Control header field can be extended through the use of one
	// or more cache-extension tokens, each with an optional value.  A cache
	// MUST ignore unrecognized cache directives.
	Extensions []string
}

func (directive *ResponseCacheDirectives) BuildResponseHeader() string {
	headerTokens := make([]string, 0)
	if directive.Immutable {
		directive.MaxAge = DeltaSeconds(defaultMaxAge.Seconds())
	}

	if len(directive.Private) > 0 {
		headerTokens = append(headerTokens, "private")
	}

	if len(directive.NoCache) > 0 {
		headerTokens = append(headerTokens, "no-cache")
	}

	if directive.Public {
		headerTokens = append(headerTokens, "public")
	}

	if directive.NoStore {
		headerTokens = append(headerTokens, "no-store")
	}

	if directive.NoTransform {
		headerTokens = append(headerTokens, "no-transform")
	}

	if directive.MustRevalidate {
		headerTokens = append(headerTokens, "must-revalidate")
	}

	if directive.ProxyRevalidate {
		headerTokens = append(headerTokens, "proxy-revalidate")
	}

	if directive.Immutable {
		headerTokens = append(headerTokens, "immutable")
	}

	if directive.MaxAge != 0 {
		headerTokens = append(headerTokens, fmt.Sprintf("max-age=%d", directive.MaxAge))
	}

	if directive.SMaxAge != 0 {
		headerTokens = append(headerTokens, fmt.Sprintf("s-maxage=%d", directive.SMaxAge))
	}

	if directive.StaleIfError != 0 {
		headerTokens = append(headerTokens, fmt.Sprintf("stale-if-error=%d", directive.StaleIfError))
	}

	if directive.StaleWhileRevalidate != 0 {
		headerTokens = append(headerTokens, fmt.Sprintf("stale-while-revalidate=%d", directive.StaleWhileRevalidate))
	}

	headerTokens = append(headerTokens, directive.Extensions...)

	return strings.Join(headerTokens, ", ")
}
